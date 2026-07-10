// Package get implements calls to the get method.
package getfeature

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
		Use:  "get-feature",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := sidecar.NewClient(sidecar.ClientOptions{Address: address, Insecure: insecure, Headers: headers})
			for j := 0; j < n; j++ {
				response, err := sidecar.CallUnary[routeguidepb.Point, routeguidepb.Feature](
					cmd.Context(),
					client,
					constants.RouteGuideGetFeatureProcedure,
					sidecar.NewRequest(&routeguidepb.Point{}),
				)
				if err != nil {
					log.Printf("returning from 1 %T %+v", err, err)
					return err
				}
				if n == 1 {
					body, err := protojson.Marshal(response.Msg)
					if err != nil {
						log.Printf("returning from 2")

						return err
					}
					_, _ = cmd.OutOrStdout().Write(body)
					_, _ = cmd.OutOrStdout().Write([]byte("\n"))
					if verbose {
						fmt.Println("Response Trailers:")
						for key, values := range response.Trailer {
							fmt.Printf("  %s: %v\n", key, values)
						}
					}
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
