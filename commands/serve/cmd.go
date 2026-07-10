// Package serve implements an Echo server.
package serve

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"charm.land/log/v2"
	"github.com/agentio/routeguide-sidecar/constants"
	routeguidepb "github.com/agentio/routeguide-sidecar/genproto/routeguidepb"
	"github.com/agentio/sidecar"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var port int
	var socket string
	var verbose bool
	cmd := &cobra.Command{
		Use:  "serve",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			mux := http.NewServeMux()
			mux.HandleFunc(constants.RouteGuideGetFeatureProcedure, sidecar.HandleUnary(getFeature))
			mux.HandleFunc(constants.RouteGuideListFeaturesProcedure, sidecar.HandleServerStreaming(listFeatures))
			mux.HandleFunc(constants.RouteGuideRecordRouteProcedure, sidecar.HandleClientStreaming(recordRoute))
			mux.HandleFunc(constants.RouteGuideRouteChatProcedure, sidecar.HandleBidiStreaming(routeChat))
			server := sidecar.NewServer(mux)
			var err error
			var listener net.Listener
			if port == 0 {
				listener, err = net.Listen("unix", socket)
			} else {
				listener, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
			}
			if err != nil {
				return err
			}
			return server.Serve(listener)
		},
	}
	cmd.Flags().IntVarP(&port, "port", "p", 0, "server port")
	cmd.Flags().StringVarP(&socket, "socket", "s", "@echo", "server socket")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose")
	return cmd
}

func getFeature(ctx context.Context, req *sidecar.Request[routeguidepb.Point]) (*sidecar.Response[routeguidepb.Feature], error) {
	log.Infof("%s", constants.RouteGuideGetFeatureProcedure)
	return sidecar.NewResponse(&routeguidepb.Feature{}), nil
}

func listFeatures(ctx context.Context, req *sidecar.Request[routeguidepb.Rectangle], stream *sidecar.ServerStream[routeguidepb.Feature]) error {
	log.Infof("%s", constants.RouteGuideListFeaturesProcedure)
	parts := strings.Split("1 2 3", " ")
	for _, part := range parts {
		_ = part
		if err := stream.Send(&routeguidepb.Feature{}); err != nil {
			return err
		}
	}
	return nil
}

func recordRoute(ctx context.Context, stream *sidecar.ClientStream[routeguidepb.Point]) (*sidecar.Response[routeguidepb.RouteSummary], error) {
	log.Infof("%s", constants.RouteGuideRecordRouteProcedure)
	parts := []string{}
	for {
		request, err := stream.Receive()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, err
		}
		_ = request
		parts = append(parts, "todo")
	}
	return sidecar.NewResponse(&routeguidepb.RouteSummary{}), nil
}

func routeChat(ctx context.Context, stream *sidecar.BidiStream[routeguidepb.RouteNote, routeguidepb.RouteNote]) error {
	log.Infof("%s", constants.RouteGuideRouteChatProcedure)
	for {
		request, err := stream.Receive()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return err
		}
		_ = request
		err = stream.Send(&routeguidepb.RouteNote{})
		if err != nil {
			return err
		}
	}
	return nil
}
