package cmd

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/mati2251/notion-to-google-tasks/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
)

var googleClient, _ = utils.GetGoogleToken()
var notionClient, _ = utils.GetNotionToken()

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration of the application",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if viper.ConfigFileUsed() != "" {
			checkOldConfigAndRemoveIt()
		}
		services, _ := cmd.Flags().GetStringSlice("specific")
		if slices.Contains(services, "google") {
			var err error
			googleClient, err = utils.GoogleConfig()
			if err != nil {
				fmt.Printf("Something went wrong: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Google config done")
		}
		if slices.Contains(services, "notion") {
			var err error
			notionClient, err = utils.NotionConfig()
			if err != nil {
				fmt.Printf("Something went wrong: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Notion config done")
		}
		if slices.Contains(services, "connections") {
			if googleClient == nil || notionClient == nil {
				fmt.Printf("Notion or google clients don't set")
				os.Exit(1)
			}
			utils.ConfigConnections(googleClient, notionClient)
		}
		viper.SafeWriteConfig()
		viper.WriteConfig()
	},
}

func init() {
	configCmd.Flags().StringSliceP("specific", "s", []string{"google", "notion", "connections"}, "Specific pages to sync")
	rootCmd.AddCommand(configCmd)
}

func checkOldConfigAndRemoveIt() {
	prompt := promptui.Prompt{
		Label:     "Do you want delete old config file?",
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
			fmt.Printf("Error on remove old config file %v\n", err)
			os.Exit(1)
		}
	}
}
