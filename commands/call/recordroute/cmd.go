// Package collect implements calls to the collect method.
package recordroute

import (
	"fmt"
	"log"

	"github.com/agentio/routeguide-sidecar/constants"
	routeguidepb "github.com/agentio/routeguide-sidecar/genproto/routeguidepb"
	"github.com/agentio/sidecar"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

func Cmd() *cobra.Command {
	var message string
	var address string
	var n int
	var verbose bool
	var insecure bool
	var headers []string
	cmd := &cobra.Command{
		Use:  "record-route",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := sidecar.NewClient(sidecar.ClientOptions{Address: address, Insecure: insecure, Headers: headers})
			stream, err := sidecar.CallClientStream[routeguidepb.Point, routeguidepb.RouteSummary](
				cmd.Context(),
				client,
				constants.RouteGuideRecordRouteProcedure,
			)
			if err != nil {
				return err
			}
			for range 3 {
				err = stream.Send(&routeguidepb.Point{})
				if err != nil {
					log.Printf("Error writing to pipe: %v", err)
					return err
				}
			}
			response, err := stream.CloseAndReceive()
			if err != nil {
				return err
			}
			body, err := protojson.Marshal(response)
			if err != nil {
				return err
			}
			_, _ = cmd.OutOrStdout().Write(body)
			_, _ = cmd.OutOrStdout().Write([]byte("\n"))
			if verbose {
				fmt.Println("Response Trailers:")
				for key, values := range stream.Trailer {
					fmt.Printf("  %s: %v\n", key, values)
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&message, "message", "m", "hello", "message")
	cmd.Flags().StringVarP(&address, "address", "a", "unix:@echo", "address of the echo server to use")
	cmd.Flags().IntVarP(&n, "number", "n", 1, "number of times to call the method")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose")
	cmd.Flags().BoolVarP(&insecure, "insecure", "i", false, "disable TLS certificate verification")
	cmd.Flags().StringArrayVarP(&headers, "header", "H", []string{}, "headers to add to the request")
	return cmd
}
