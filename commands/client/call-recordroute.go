package client

import (
	"context"
	"log"
	rand "math/rand/v2"
	"time"

	routeguidepb "github.com/agentio/routeguide-sidecar/genproto/routeguidepb"
	"github.com/agentio/routeguide-sidecar/internal/routeguide/constants"
	"github.com/agentio/sidecar"
)

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
	// call RouteGuide/RecordRoute
	stream, err := sidecar.CallClientStream[routeguidepb.Point, routeguidepb.RouteSummary](
		ctx, client, constants.RouteGuideRecordRouteProcedure,
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

func randomPoint() *routeguidepb.Point {
	lat := (rand.Int32N(180) - 90) * 1e7
	long := (rand.Int32N(360) - 180) * 1e7
	return &routeguidepb.Point{Latitude: lat, Longitude: long}
}
