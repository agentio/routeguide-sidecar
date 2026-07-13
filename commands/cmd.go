// Package commands contains command implementations.
package commands

import (
	"github.com/agentio/routeguide-sidecar/commands/client"
	"github.com/agentio/routeguide-sidecar/commands/serve"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "routeguide-sidecar",
	}
	cmd.AddCommand(client.Cmd())
	cmd.AddCommand(serve.Cmd())
	return cmd
}
