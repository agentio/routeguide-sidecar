package client

import (
	"context"
	"io"
	"log"
	rand "math/rand/v2"
	"time"

	"github.com/agentio/routeguide-sidecar/constants"
	routeguidepb "github.com/agentio/routeguide-sidecar/genproto/routeguidepb"
	"github.com/agentio/sidecar"
	"github.com/spf13/cobra"
)

// printFeature gets the feature for the given point.
func printFeature(client *sidecar.Client, point *routeguidepb.Point) {
	log.Printf("Getting feature for point (%d, %d)", point.Latitude, point.Longitude)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//feature, err := client.GetFeature(ctx, point)
	response, err := sidecar.CallUnary[routeguidepb.Point, routeguidepb.Feature](
		ctx,
		client,
		constants.RouteGuideGetFeatureProcedure,
		sidecar.NewRequest(point),
	)
	if err != nil {
		log.Fatalf("client.GetFeature failed: %v", err)
	}
	log.Println(response.Msg)
}

// printFeatures lists all the features within the given bounding Rectangle.
func printFeatures(client *sidecar.Client, rect *routeguidepb.Rectangle) {
	log.Printf("Looking for features within %v", rect)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := sidecar.CallServerStream[routeguidepb.Rectangle, routeguidepb.Feature](
		ctx,
		client,
		constants.RouteGuideListFeaturesProcedure,
		sidecar.NewRequest(rect),
	)
	if err != nil {
		log.Fatalf("client.ListFeatures failed: %v", err)
	}
	for {
		feature, err := stream.Receive()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("client.ListFeatures failed: %v", err)
		}
		log.Printf("Feature: name: %q, point:(%v, %v)", feature.GetName(),
			feature.GetLocation().GetLatitude(), feature.GetLocation().GetLongitude())
	}
}

// runRecordRoute sends a sequence of points to server and expects to get a RouteSummary from server.
func runRecordRoute(client *sidecar.Client) {
	// Create a random number of random points
	pointCount := int(rand.Int32N(100)) + 2 // Traverse at least two points
	var points []*routeguidepb.Point
	for i := 0; i < pointCount; i++ {
		points = append(points, randomPoint())
	}
	log.Printf("Traversing %d points.", len(points))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := sidecar.CallClientStream[routeguidepb.Point, routeguidepb.RouteSummary](
		ctx,
		client,
		constants.RouteGuideRecordRouteProcedure,
	)
	if err != nil {
		log.Fatalf("client.RecordRoute failed: %v", err)
	}
	for _, point := range points {
		if err := stream.Send(point); err != nil {
			log.Fatalf("client.RecordRoute: stream.Send(%v) failed: %v", point, err)
		}
	}
	reply, err := stream.CloseAndReceive()
	if err != nil {
		log.Fatalf("client.RecordRoute failed: %v", err)
	}
	log.Printf("Route summary: %v", reply)
}

// runRouteChat receives a sequence of route notes, while sending notes for various locations.
func runRouteChat(client *sidecar.Client) {
	notes := []*routeguidepb.RouteNote{
		{Location: &routeguidepb.Point{Latitude: 0, Longitude: 1}, Message: "First message"},
		{Location: &routeguidepb.Point{Latitude: 0, Longitude: 2}, Message: "Second message"},
		{Location: &routeguidepb.Point{Latitude: 0, Longitude: 3}, Message: "Third message"},
		{Location: &routeguidepb.Point{Latitude: 0, Longitude: 1}, Message: "Fourth message"},
		{Location: &routeguidepb.Point{Latitude: 0, Longitude: 2}, Message: "Fifth message"},
		{Location: &routeguidepb.Point{Latitude: 0, Longitude: 3}, Message: "Sixth message"},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := sidecar.CallBidiStream[routeguidepb.RouteNote, routeguidepb.RouteNote](
		ctx,
		client,
		constants.RouteGuideRouteChatProcedure,
	)
	if err != nil {
		log.Fatalf("client.RouteChat failed: %v", err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Receive()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("client.RouteChat failed: %v", err)
			}
			log.Printf("Got message %s at point(%d, %d)", in.Message, in.Location.Latitude, in.Location.Longitude)
		}
	}()
	for _, note := range notes {
		if err := stream.Send(note); err != nil {
			log.Fatalf("client.RouteChat: stream.Send(%v) failed: %v", note, err)
		}
	}
	stream.CloseRequest()
	<-waitc
}

func randomPoint() *routeguidepb.Point {
	lat := (rand.Int32N(180) - 90) * 1e7
	long := (rand.Int32N(360) - 180) * 1e7
	return &routeguidepb.Point{Latitude: lat, Longitude: long}
}

func Cmd() *cobra.Command {
	var message string
	var address string
	var insecure bool
	var headers []string
	cmd := &cobra.Command{
		Use:  "client",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := sidecar.NewClient(sidecar.ClientOptions{Address: address, Insecure: insecure, Headers: headers})
			// Looking for a valid feature
			printFeature(client, &routeguidepb.Point{Latitude: 409146138, Longitude: -746188906})
			// Feature missing.
			printFeature(client, &routeguidepb.Point{Latitude: 0, Longitude: 0})
			// Looking for features between 40, -75 and 42, -73.
			printFeatures(client, &routeguidepb.Rectangle{
				Lo: &routeguidepb.Point{Latitude: 400000000, Longitude: -750000000},
				Hi: &routeguidepb.Point{Latitude: 420000000, Longitude: -730000000},
			})
			// RecordRoute
			runRecordRoute(client)
			// RouteChat
			runRouteChat(client)
			return nil
		},
	}
	cmd.Flags().StringVarP(&message, "message", "m", "hello", "message")
	cmd.Flags().StringVarP(&address, "address", "a", "unix:@echo", "address of the echo server to use")
	cmd.Flags().BoolVarP(&insecure, "insecure", "i", false, "disable TLS certificate verification")
	cmd.Flags().StringArrayVarP(&headers, "header", "H", []string{}, "headers to add to the request")
	return cmd
}
