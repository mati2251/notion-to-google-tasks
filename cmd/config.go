package cmd

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/mati2251/notion-to-google-tasks/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration of the application",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if viper.ConfigFileUsed() != "" {
			checkOldConfigAndRemoveIt()
		}
		utils.GoogleConfig()
		viper.SafeWriteConfig()
		viper.WriteConfig()
	},
}

func init() {
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
		os.Exit(0)
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
