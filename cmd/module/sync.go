package module

import (
	"fmt"

	"github.com/spf13/cobra"
)

// SyncCmd returns the cobra command that reconciles the DB with on-disk manifests.
func SyncCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "Reconcile database with on-disk module manifests",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Syncing module registry with on-disk manifests...")
			return nil
		},
	}
}
