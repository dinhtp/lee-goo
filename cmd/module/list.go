package module

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ListCmd returns the cobra command that lists all discovered modules.
func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all discovered modules",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Listing modules...")
			return nil
		},
	}
}
