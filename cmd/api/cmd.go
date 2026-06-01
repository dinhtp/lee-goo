package api

import "github.com/spf13/cobra"

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api",
		Short: "Manage the API server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(ServeCmd())
	return cmd
}
