package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/mati2251/notion-to-google-tasks/config/auth"
	"github.com/mati2251/notion-to-google-tasks/config/connections"
	"github.com/mati2251/notion-to-google-tasks/sync"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration of the application",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		auth.InitConnections()
		specifics, _ := cmd.Flags().GetStringSlice("specific")
		removes, _ := cmd.Flags().GetStringSlice("remove")
		if slices.Contains(removes, "all") && viper.ConfigFileUsed() != "" {
			checkOldConfigAndRemoveIt()
		} else {
			if slices.Contains(removes, "google") {
				auth.RemoveGoogleConfig()
				fmt.Println("Google config removed")
			}
			if slices.Contains(removes, "notion") {
				auth.RemoveNotionConfig()
				fmt.Println("Notion config removed")
			}
			if slices.Contains(removes, "connections") {
				connections.RemoveConnections()
				fmt.Println("Connections config removed")
			}
			viper.WriteConfig()
		}
		if slices.Contains(specifics, "google") {
			auth.GoogleConfig()
			fmt.Println("Google config done")
		}
		if slices.Contains(specifics, "notion") {
			auth.NotionConfig()
			fmt.Println("Notion config done")
		}
		if slices.Contains(specifics, "connections") {
			if auth.NotionClient == nil || auth.TasksService == nil {
				log.Fatalf("Notion or google clients don't set")
			}
			connections.ConfigConnections()
		}
		if slices.Contains(specifics, "first-scan") {
			sync.ForceSync()
		}
	},
}

func init() {
	configCmd.Flags().StringSliceP(
		"specific",
		"s",
		[]string{"google", "notion", "connections", "first-scan"},
		"Specific pages to sync (avaliable:google,notion,connections,first-scan,none)",
	)
	configCmd.Flags().StringSliceP("remove", "r", []string{"all"}, "Remove old config specific part (avaliable:all,google,notion,connections)")
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
		viper.Reset()
		viper.SetConfigFile(fileName)
	}
}
