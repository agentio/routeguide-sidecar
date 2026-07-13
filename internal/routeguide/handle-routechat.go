package routeguide

import (
	"context"
	"fmt"
	"io"
	"log"

	routeguidepb "github.com/agentio/routeguide-sidecar/genproto/routeguidepb"
	"github.com/agentio/routeguide-sidecar/internal/routeguide/constants"
	"github.com/agentio/sidecar"
)

// RouteChat receives a stream of message/location pairs, and responds with a stream of all
// previous messages at each of those locations.
func (s *Server) routeChat(ctx context.Context, stream *sidecar.BidiStream[routeguidepb.RouteNote, routeguidepb.RouteNote]) error {
	log.Printf("%s", constants.RouteGuideRouteChatProcedure)
	for {
		in, err := stream.Receive()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		key := serialize(in.Location)
		s.mu.Lock()
		s.routeNotes[key] = append(s.routeNotes[key], in)
		// Note: this copy prevents blocking other clients while serving this one.
		// We don't need to do a deep copy, because elements in the slice are
		// insert-only and never modified.
		rn := make([]*routeguidepb.RouteNote, len(s.routeNotes[key]))
		copy(rn, s.routeNotes[key])
		s.mu.Unlock()
		for _, note := range rn {
			if err := stream.Send(note); err != nil {
				return err
			}
		}
	}
}

func serialize(point *routeguidepb.Point) string {
	return fmt.Sprintf("%d %d", point.Latitude, point.Longitude)
}
