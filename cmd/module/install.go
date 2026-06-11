package module

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// InstallCmd returns the cobra command that installs a module by running its DB migrations.
func InstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install <name>",
		Short: "Install a module (runs its DB migrations)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			runner, err := buildRunner(name)
			if err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			if err := runner.UpFor(ctx, name); err != nil {
				return fmt.Errorf("install %s: %w", name, err)
			}
			fmt.Printf("Module %q installed successfully.\n", name)
			return nil
		},
	}
}
