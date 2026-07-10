// routeguide-sidecar is a command-line tool that uses sidecar to implement and call a gRPC service.
package main

import (
	"os"

	"github.com/agentio/routeguide-sidecar/commands"
)

func main() {
	if err := commands.Cmd().Execute(); err != nil {
		os.Exit(1)
	}
}
