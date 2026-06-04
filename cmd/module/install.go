package module

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/fx"

	authzModule "github.com/dinhtp/lee-goo/modules/authorization/fx"
	coreModule "github.com/dinhtp/lee-goo/modules/core/fx"
	userModule "github.com/dinhtp/lee-goo/modules/user/fx"
	systemfx "github.com/dinhtp/lee-goo/system/fx"
	"github.com/dinhtp/lee-goo/system/migrator"
)

// InstallCmd returns the cobra command that installs a module by running its DB migrations.
func InstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install <name>",
		Short: "Install a module (runs its DB migrations)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			var runErr error
			app := fx.New(
				systemfx.MigrateOptions(),
				coreModule.MigrationSource(),
				userModule.MigrationSource(),
				authzModule.MigrationSource(),
				fx.NopLogger,
				fx.Invoke(func(runner *migrator.Runner) {
					runErr = runner.UpFor(cmd.Context(), name)
				}),
			)
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			if err := app.Start(ctx); err != nil {
				return fmt.Errorf("install %s: %w", name, err)
			}
			_ = app.Stop(ctx)
			if runErr != nil {
				return fmt.Errorf("install %s: migration failed: %w", name, runErr)
			}
			fmt.Printf("Module %q installed successfully.\n", name)
			return nil
		},
	}
}
