package module

import (
	"fmt"

	"github.com/spf13/cobra"
)

// MigrateCmd returns the cobra command that runs database migrations for a module.
func MigrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate [name]",
		Short: "Run database migrations for a module (or all modules)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				fmt.Println("Running migrations for all modules...")
			} else {
				fmt.Printf("Running migrations for module %q...\n", args[0])
			}
			return nil
		},
	}
	return cmd
}
