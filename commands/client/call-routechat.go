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
	// call RouteGuide/RouteChat
	stream, err := sidecar.CallBidiStream[routeguidepb.RouteNote, routeguidepb.RouteNote](
		ctx, client, constants.RouteGuideRouteChatProcedure,
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
