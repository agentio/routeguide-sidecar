// Package call implements calls to the Echo service.
package call

import (
	"github.com/agentio/routeguide-sidecar/commands/call/getfeature"
	"github.com/agentio/routeguide-sidecar/commands/call/listfeatures"
	"github.com/agentio/routeguide-sidecar/commands/call/recordroute"
	"github.com/agentio/routeguide-sidecar/commands/call/routechat"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "call",
	}
	cmd.AddCommand(getfeature.Cmd())
	cmd.AddCommand(listfeatures.Cmd())
	cmd.AddCommand(recordroute.Cmd())
	cmd.AddCommand(routechat.Cmd())
	return cmd
}
