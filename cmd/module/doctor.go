package module

import (
	"fmt"

	"github.com/spf13/cobra"
)

// DoctorCmd returns the cobra command that checks for module health issues.
func DoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check module system for issues",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Running module health checks...")
			return nil
		},
	}
}
