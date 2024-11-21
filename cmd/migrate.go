package cmd

import (
	"github.com/spf13/cobra"

	"github.com/fmiskovic/cash-me-if-you-can/cmd/migrate"
)

func init() {
	migrateCmd.AddCommand(migrate.UpCmd)
	migrateCmd.AddCommand(migrate.DownCmd)

	rootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate commands",
	Long:  `Migrate commands run the migrations`,
}
