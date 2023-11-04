package cmd

import (
	"github.com/mati2251/notion-to-google-tasks/config/auth"
	"github.com/mati2251/notion-to-google-tasks/config/connections"
	"github.com/mati2251/notion-to-google-tasks/sync"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronize notion databases with google tasks. Requiers valid config(see config command)",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := auth.InitConnections()
		cobra.CheckErr(err)
		connections := connections.GetConnections()
		err = sync.Sync(connections)
		cobra.CheckErr(err)
	},
}

func init() {
	syncCmd.Flags().BoolP("force", "f", false, "Force sync even if it is not needed")
	rootCmd.AddCommand(syncCmd)
}
