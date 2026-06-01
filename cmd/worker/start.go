package worker

import (
	"fmt"

	"github.com/spf13/cobra"
)

func StartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start the background worker (not yet implemented)",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: implement worker startup logic
			fmt.Println("worker: not implemented")
			return nil
		},
	}
}
