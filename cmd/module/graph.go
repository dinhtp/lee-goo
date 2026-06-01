package module

import (
	"fmt"

	"github.com/spf13/cobra"
)

// GraphCmd returns the cobra command that prints the module dependency graph.
func GraphCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "graph",
		Short: "Print the module dependency graph",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Building dependency graph...")
			return nil
		},
	}
}
