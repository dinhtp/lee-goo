package module

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// UninstallCmd returns the cobra command that uninstalls a module by rolling back its DB migrations.
func UninstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "uninstall <name>",
		Short: "Uninstall a module (rolls back its DB migrations)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			runner, err := buildRunner(name)
			if err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			if err := runner.DownFor(ctx, name); err != nil {
				return fmt.Errorf("uninstall %s: %w", name, err)
			}
			fmt.Printf("Module %q uninstalled successfully.\n", name)
			return nil
		},
	}
}
