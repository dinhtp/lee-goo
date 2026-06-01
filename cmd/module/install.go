package module

import (
	"fmt"

	"github.com/spf13/cobra"
)

// InstallCmd returns the cobra command that installs a module.
func InstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install <name>",
		Short: "Install a module",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Installing module %q...\n", args[0])
			return nil
		},
	}
}
