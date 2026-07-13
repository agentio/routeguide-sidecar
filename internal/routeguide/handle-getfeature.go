package routeguide

import (
	"context"
	"log"

	routeguidepb "github.com/agentio/routeguide-sidecar/genproto/routeguidepb"
	"github.com/agentio/routeguide-sidecar/internal/routeguide/constants"
	"github.com/agentio/sidecar"
	"google.golang.org/protobuf/proto"
)

// GetFeature returns the feature at the given point.
func (s *Server) getFeature(ctx context.Context, req *sidecar.Request[routeguidepb.Point]) (*sidecar.Response[routeguidepb.Feature], error) {
	log.Printf("%s", constants.RouteGuideGetFeatureProcedure)
	for _, feature := range s.savedFeatures {
		if proto.Equal(feature.Location, req.Msg) {
			return sidecar.NewResponse(feature), nil
		}
	}
	// No feature was found, return an unnamed feature
	return sidecar.NewResponse(&routeguidepb.Feature{Location: req.Msg}), nil
}
