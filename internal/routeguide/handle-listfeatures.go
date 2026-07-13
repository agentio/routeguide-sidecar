package routeguide

import (
	"context"
	"log"
	"math"

	routeguidepb "github.com/agentio/routeguide-sidecar/genproto/routeguidepb"
	"github.com/agentio/routeguide-sidecar/internal/routeguide/constants"
	"github.com/agentio/sidecar"
)

// ListFeatures lists all features contained within the given bounding Rectangle.
func (s *Server) listFeatures(ctx context.Context, req *sidecar.Request[routeguidepb.Rectangle], stream *sidecar.ServerStream[routeguidepb.Feature]) error {
	log.Printf("%s", constants.RouteGuideListFeaturesProcedure)
	for _, feature := range s.savedFeatures {
		if inRange(feature.Location, req.Msg) {
			if err := stream.Send(feature); err != nil {
				return err
			}
		}
	}
	return nil
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
