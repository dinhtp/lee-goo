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
		StatusCmd(),
		InstallCmd(),
		EnableCmd(),
		DisableCmd(),
		UninstallCmd(),
		MigrateCmd(),
		GraphCmd(),
		DoctorCmd(),
		MakeCmd(),
		SyncCmd(),
	)
	return cmd
}
