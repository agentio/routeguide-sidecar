package client

import (
	"context"
	"log"
	"time"

	routeguidepb "github.com/agentio/routeguide-sidecar/genproto/routeguidepb"
	"github.com/agentio/routeguide-sidecar/internal/routeguide/constants"
	"github.com/agentio/sidecar"
)

// printFeature gets the feature for the given point.
func printFeature(client *sidecar.Client, point *routeguidepb.Point) {
	log.Printf("Getting feature for point (%d, %d)", point.Latitude, point.Longitude)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// call RouteGuide/GetFeature
	response, err := sidecar.CallUnary[routeguidepb.Point, routeguidepb.Feature](
		ctx, client, constants.RouteGuideGetFeatureProcedure,
		sidecar.NewRequest(point),
	)
	if err != nil {
		log.Fatalf("client.GetFeature failed: %s", err)
	}
	log.Println(response.Msg)
}
