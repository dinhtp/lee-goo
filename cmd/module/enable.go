package module

import (
	"fmt"

	"github.com/spf13/cobra"
)

// EnableCmd returns the cobra command that enables a module.
func EnableCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "enable <name>",
		Short: "Enable an installed module",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Enabling module %q...\n", args[0])
			return nil
		},
	}
}
