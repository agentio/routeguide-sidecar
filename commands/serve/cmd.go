// Package serve implements a Route Guide server.
package serve

import (
	"fmt"
	"net"

	"github.com/agentio/routeguide-sidecar/internal/routeguide"
	"github.com/agentio/sidecar"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var port int
	var socket string
	var verbose bool
	var data string
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Run the routeguide server",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			apiServer, err := routeguide.NewServer(data)
			if err != nil {
				return err
			}
			httpServer := sidecar.NewServer(apiServer.ServeMux())
			var listener net.Listener
			if port == 0 {
				listener, err = net.Listen("unix", socket)
			} else {
				listener, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
			}
			if err != nil {
				return err
			}
			return httpServer.Serve(listener)
		},
	}
	cmd.Flags().IntVarP(&port, "port", "p", 0, "server port")
	cmd.Flags().StringVarP(&socket, "socket", "s", "@routeguide", "server socket")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose")
	cmd.Flags().StringVar(&data, "data", "", "path to a JSON file containing point")
	return cmd
}
