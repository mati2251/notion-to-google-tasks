package cmd

import (
	"log"

	"github.com/mati2251/notion-to-google-tasks/utils"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronize notion databases with google tasks. Requiers valid config(see config command)",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Init()
		isForce, err := cmd.Flags().GetBool("force")
		if err != nil {
			log.Fatalf("Error getting force flag: %v", err)
		}
		if isForce {
			utils.ForceSync()
		} else {
			utils.Sync()
		}
	},
}

func init() {
	syncCmd.Flags().BoolP("force", "f", false, "Force sync even if it is not needed")
	rootCmd.AddCommand(syncCmd)
}
