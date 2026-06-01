package module

import (
	"fmt"

	"github.com/spf13/cobra"
)

// StatusCmd returns the cobra command that shows a module's current status.
func StatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status <name>",
		Short: "Show status of a module",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Status of module %q...\n", args[0])
			return nil
		},
	}
}
