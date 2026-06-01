package module

import (
	"fmt"

	"github.com/spf13/cobra"
)

// DisableCmd returns the cobra command that disables an enabled module.
func DisableCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "disable <name>",
		Short: "Disable an enabled module",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Disabling module %q...\n", args[0])
			return nil
		},
	}
}
