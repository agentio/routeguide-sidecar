// Package serve implements an Echo server.
package serve

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"charm.land/log/v2"
	"github.com/agentio/routeguide-sidecar/commands/serve/testdata"
	"github.com/agentio/routeguide-sidecar/constants"
	routeguidepb "github.com/agentio/routeguide-sidecar/genproto/routeguidepb"
	"github.com/agentio/sidecar"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
)

func Cmd() *cobra.Command {
	var port int
	var socket string
	var verbose bool
	cmd := &cobra.Command{
		Use:  "serve",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			mux := http.NewServeMux()
			mux.HandleFunc(constants.RouteGuideGetFeatureProcedure, sidecar.HandleUnary(getFeature))
			mux.HandleFunc(constants.RouteGuideListFeaturesProcedure, sidecar.HandleServerStreaming(listFeatures))
			mux.HandleFunc(constants.RouteGuideRecordRouteProcedure, sidecar.HandleClientStreaming(recordRoute))
			mux.HandleFunc(constants.RouteGuideRouteChatProcedure, sidecar.HandleBidiStreaming(routeChat))
			server := sidecar.NewServer(mux)
			var err error
			var listener net.Listener
			if port == 0 {
				listener, err = net.Listen("unix", socket)
			} else {
				listener, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
			}
			if err != nil {
				return err
			}
			return server.Serve(listener)
		},
	}
	cmd.Flags().IntVarP(&port, "port", "p", 0, "server port")
	cmd.Flags().StringVarP(&socket, "socket", "s", "@routeguide", "server socket")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose")
	return cmd
}

var savedFeatures []*routeguidepb.Feature
var mu sync.Mutex // protects routeNotes
var routeNotes map[string][]*routeguidepb.RouteNote

func init() {
	savedFeatures = []*routeguidepb.Feature{}
	routeNotes = map[string][]*routeguidepb.RouteNote{}
	loadFeatures("")
}

func getFeature(ctx context.Context, req *sidecar.Request[routeguidepb.Point]) (*sidecar.Response[routeguidepb.Feature], error) {
	log.Infof("%s", constants.RouteGuideGetFeatureProcedure)
	point := req.Msg
	for _, feature := range savedFeatures {
		if proto.Equal(feature.Location, point) {
			return sidecar.NewResponse(feature), nil
		}
	}
	// No feature was found, return an unnamed feature
	return sidecar.NewResponse(&routeguidepb.Feature{Location: point}), nil
}

func listFeatures(ctx context.Context, req *sidecar.Request[routeguidepb.Rectangle], stream *sidecar.ServerStream[routeguidepb.Feature]) error {
	log.Infof("%s", constants.RouteGuideListFeaturesProcedure)
	rect := req.Msg
	for _, feature := range savedFeatures {
		if inRange(feature.Location, rect) {
			if err := stream.Send(feature); err != nil {
				return err
			}
		}
	}
	return nil
}

// RecordRoute records a route composited of a sequence of points.
//
// It gets a stream of points, and responds with statistics about the "trip":
// number of points,  number of known features visited, total distance traveled, and
// total time spent.
func recordRoute(ctx context.Context, stream *sidecar.ClientStream[routeguidepb.Point]) (*sidecar.Response[routeguidepb.RouteSummary], error) {
	log.Infof("%s", constants.RouteGuideRecordRouteProcedure)
	var pointCount, featureCount, distance int32
	var lastPoint *routeguidepb.Point
	startTime := time.Now()
	for {
		point, err := stream.Receive()
		if err == io.EOF {
			endTime := time.Now()
			return sidecar.NewResponse(&routeguidepb.RouteSummary{
				PointCount:   pointCount,
				FeatureCount: featureCount,
				Distance:     distance,
				ElapsedTime:  int32(endTime.Sub(startTime).Seconds()),
			}), nil
		}
		if err != nil {
			return nil, err
		}
		pointCount++
		for _, feature := range savedFeatures {
			if proto.Equal(feature.Location, point) {
				featureCount++
			}
		}
		if lastPoint != nil {
			distance += calcDistance(lastPoint, point)
		}
		lastPoint = point
	}
}

// RouteChat receives a stream of message/location pairs, and responds with a stream of all
// previous messages at each of those locations.
func routeChat(ctx context.Context, stream *sidecar.BidiStream[routeguidepb.RouteNote, routeguidepb.RouteNote]) error {
	log.Infof("%s", constants.RouteGuideRouteChatProcedure)
	for {
		in, err := stream.Receive()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		key := serialize(in.Location)

		mu.Lock()
		routeNotes[key] = append(routeNotes[key], in)
		// Note: this copy prevents blocking other clients while serving this one.
		// We don't need to do a deep copy, because elements in the slice are
		// insert-only and never modified.
		rn := make([]*routeguidepb.RouteNote, len(routeNotes[key]))
		copy(rn, routeNotes[key])
		mu.Unlock()

		for _, note := range rn {
			if err := stream.Send(note); err != nil {
				return err
			}
		}
	}
}

// loadFeatures loads features from a JSON file.
func loadFeatures(filePath string) {
	var data []byte
	if filePath != "" {
		var err error
		data, err = os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Failed to load default features: %v", err)
		}
	} else {
		data = testdata.Data
	}
	if err := json.Unmarshal(data, &savedFeatures); err != nil {
		log.Fatalf("Failed to load default features: %v", err)
	}
}

func toRadians(num float64) float64 {
	return num * math.Pi / float64(180)
}

// calcDistance calculates the distance between two points using the "haversine" formula.
// The formula is based on http://mathforum.org/library/drmath/view/51879.html.
func calcDistance(p1 *routeguidepb.Point, p2 *routeguidepb.Point) int32 {
	const CordFactor float64 = 1e7
	const R = float64(6371000) // earth radius in metres
	lat1 := toRadians(float64(p1.Latitude) / CordFactor)
	lat2 := toRadians(float64(p2.Latitude) / CordFactor)
	lng1 := toRadians(float64(p1.Longitude) / CordFactor)
	lng2 := toRadians(float64(p2.Longitude) / CordFactor)
	dlat := lat2 - lat1
	dlng := lng2 - lng1

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dlng/2)*math.Sin(dlng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := R * c
	return int32(distance)
}

func inRange(point *routeguidepb.Point, rect *routeguidepb.Rectangle) bool {
	left := math.Min(float64(rect.Lo.Longitude), float64(rect.Hi.Longitude))
	right := math.Max(float64(rect.Lo.Longitude), float64(rect.Hi.Longitude))
	top := math.Max(float64(rect.Lo.Latitude), float64(rect.Hi.Latitude))
	bottom := math.Min(float64(rect.Lo.Latitude), float64(rect.Hi.Latitude))

	if float64(point.Longitude) >= left &&
		float64(point.Longitude) <= right &&
		float64(point.Latitude) >= bottom &&
		float64(point.Latitude) <= top {
		return true
	}
	return false
}

func serialize(point *routeguidepb.Point) string {
	return fmt.Sprintf("%d %d", point.Latitude, point.Longitude)
}

// exampleData is a copy of testdata/route_guide_db.json. It's to avoid
// specifying file path with `go run`.
var exampleData = []byte(``)
