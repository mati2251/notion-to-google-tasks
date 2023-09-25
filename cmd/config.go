package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/mati2251/notion-to-google-tasks/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration of the application",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var tasksService, _ = utils.GetTasksService()
		var notionClient, _ = utils.GetNotionToken()
		services, _ := cmd.Flags().GetStringSlice("specific")
		if viper.ConfigFileUsed() != "" && slices.Contains(services, "google") {
			checkOldConfigAndRemoveIt()
		}
		if slices.Contains(services, "google") {
			var err error
			tasksService, err = utils.GoogleConfig()
			if err != nil {
				log.Fatalf("Something went wrong: %v\n", err)
			}
			fmt.Println("Google config done")
		}
		if slices.Contains(services, "notion") {
			var err error
			notionClient, err = utils.NotionConfig()
			if err != nil {
				log.Fatalf("Something went wrong: %v\n", err)
			}
			fmt.Println("Notion config done")
		}
		if slices.Contains(services, "connections") {
			if tasksService == nil || notionClient == nil {
				log.Fatalf("Notion or google clients don't set")
			}
			utils.ConfigConnections(tasksService, notionClient)
		}
		viper.SafeWriteConfig()
		viper.WriteConfig()
	},
}

func init() {
	configCmd.Flags().StringSliceP("specific", "s", []string{"file-remove", "google", "notion", "connections"}, "Specific pages to sync")
	rootCmd.AddCommand(configCmd)
}

func checkOldConfigAndRemoveIt() {
	prompt := promptui.Prompt{
		Label:     "Do you want delete old config file",
		IsConfirm: true,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Old config will be keep %v\n", err)
	}

	if result == "y" {
		fileName := viper.ConfigFileUsed()
		err := os.Remove(fileName)
		if err != nil {
			log.Fatalf("Error on remove old config file %v\n", err)
		}
	}
}
