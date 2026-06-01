package worker

import "github.com/spf13/cobra"

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "worker",
		Short: "Manage background workers",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(StartCmd())
	return cmd
}
