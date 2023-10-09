package cmd

import (
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronize notion databases with google tasks. Requiers valid config(see config command)",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	syncCmd.Flags().BoolP("force", "f", false, "Force sync even if it is not needed")
	rootCmd.AddCommand(syncCmd)
}
