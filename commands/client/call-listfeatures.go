package client

import (
	"context"
	"io"
	"log"
	"time"

	routeguidepb "github.com/agentio/routeguide-sidecar/genproto/routeguidepb"
	"github.com/agentio/routeguide-sidecar/internal/routeguide/constants"
	"github.com/agentio/sidecar"
)

// printFeatures lists all the features within the given bounding Rectangle.
func printFeatures(client *sidecar.Client, rect *routeguidepb.Rectangle) {
	log.Printf("Looking for features within %v", rect)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// call RouteGuide/ListFeatures
	stream, err := sidecar.CallServerStream[routeguidepb.Rectangle, routeguidepb.Feature](
		ctx, client, constants.RouteGuideListFeaturesProcedure,
		sidecar.NewRequest(rect),
	)
	if err != nil {
		log.Fatalf("client.ListFeatures failed: %s", err)
	}
	for {
		feature, err := stream.Receive()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("client.ListFeatures failed: %s", err)
		}
		log.Printf("Feature: name: %q, point:(%v, %v)", feature.GetName(),
			feature.GetLocation().GetLatitude(), feature.GetLocation().GetLongitude())
	}
}
