package module

import "github.com/spf13/cobra"

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "module",
		Short: "Manage modules",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(
		ListCmd(),
		MakeCmd(),
		InstallCmd(),
		UninstallCmd(),
	)
	return cmd
}
