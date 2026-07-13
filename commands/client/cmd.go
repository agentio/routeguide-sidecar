package client

import (
	routeguidepb "github.com/agentio/routeguide-sidecar/genproto/routeguidepb"
	"github.com/agentio/sidecar"
	"github.com/spf13/cobra"
)

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
			// Unary: Looking for a valid feature
			printFeature(client, &routeguidepb.Point{Latitude: 409146138, Longitude: -746188906})
			// Unary: Feature missing.
			printFeature(client, &routeguidepb.Point{Latitude: 0, Longitude: 0})
			// Server Streaming: Looking for features between 40, -75 and 42, -73.
			printFeatures(client, &routeguidepb.Rectangle{
				Lo: &routeguidepb.Point{Latitude: 400000000, Longitude: -750000000},
				Hi: &routeguidepb.Point{Latitude: 420000000, Longitude: -730000000},
			})
			// Client Streaming: RecordRoute
			runRecordRoute(client)
			// Bidi Streaming: RouteChat
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
