package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/dinhtp/lee-goo/cmd/api"
	"github.com/dinhtp/lee-goo/cmd/module"
	"github.com/dinhtp/lee-goo/cmd/worker"
)

var rootCmd = &cobra.Command{
	Use:   "lee-goo",
	Short: "lee-goo — Go modular monorepo CLI",
	Long:  "lee-goo manages the API server, module lifecycle, and background workers.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() {
	rootCmd.AddCommand(
		api.Cmd(),
		module.Cmd(),
		worker.Cmd(),
	)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
