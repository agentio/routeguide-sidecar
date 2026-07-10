// Package commands contains command implementations.
package commands

import (
	"charm.land/log/v2"
	"github.com/agentio/routeguide-sidecar/commands/call"
	"github.com/agentio/routeguide-sidecar/commands/serve"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var logLevel string
	cmd := &cobra.Command{
		Use: "routeguide-sidecar",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			ll, err := log.ParseLevel(logLevel)
			if err != nil {
				return err
			}
			log.SetLevel(ll)
			return nil
		},
	}
	cmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "set log level (debug, info, warn, error, fatal)")
	cmd.AddCommand(call.Cmd())
	cmd.AddCommand(serve.Cmd())
	return cmd
}
