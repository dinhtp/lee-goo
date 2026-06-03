package api

import (
	"github.com/spf13/cobra"
	"go.uber.org/fx"

	authnModule "github.com/dinhtp/lee-goo/modules/authentication/fx"
	authzModule "github.com/dinhtp/lee-goo/modules/authorization/fx"
	moduleModule "github.com/dinhtp/lee-goo/modules/core/fx"
	userModule "github.com/dinhtp/lee-goo/modules/user/fx"
	systemfx "github.com/dinhtp/lee-goo/system/fx"
)

func ServeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Start the HTTP API server",
		RunE: func(cmd *cobra.Command, args []string) error {
			app := fx.New(
				systemfx.Options(),
				moduleModule.Module(),
				userModule.Module(),
				authnModule.Module(),
				authzModule.Module(),
			)
			app.Run()
			return nil
		},
	}
}
