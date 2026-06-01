package module

import (
	"fmt"

	"github.com/spf13/cobra"
)

// MakeCmd returns the cobra command that scaffolds a new module skeleton.
func MakeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "make <name>",
		Short: "Scaffold a new module skeleton",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Scaffolding new module %q...\n", args[0])
			return nil
		},
	}
}
