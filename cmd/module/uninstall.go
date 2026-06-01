package module

import (
	"fmt"

	"github.com/spf13/cobra"
)

// UninstallCmd returns the cobra command that uninstalls a module.
func UninstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall <name>",
		Short: "Uninstall a module",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			force, _ := cmd.Flags().GetBool("force")
			fmt.Printf("Uninstalling module %q (force=%v)...\n", args[0], force)
			return nil
		},
	}
	cmd.Flags().Bool("force", false, "Force uninstall even if dependents exist")
	return cmd
}
