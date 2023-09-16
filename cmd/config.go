package cmd

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration of the application",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if viper.ConfigFileUsed() != "" {
			err := checkOldConfigAndRemoveIt()
			if err != nil {
				os.Exit(1)
			}
		}
		setDefaults()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func setDefaults() {
	viper.SetDefault("google.auth_uri", "https://accounts.google.com/o/oauth2/auth")
	viper.SetDefault("google.token_uri", "https://oauth2.googleapis.com/token")
	viper.SetDefault("google.auth_provider_x509_cert_url", "https://www.googleapis.com/oauth2/v1/certs")
	viper.SetDefault("google.redirect_uris", "http://localhost")
	viper.WriteConfig()
	viper.SafeWriteConfig()
}

func checkOldConfigAndRemoveIt() error {
	prompt := promptui.Prompt{
		Label:     "Do you want delete old config file?",
		IsConfirm: true,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Old config will be keep %v\n", err)
		return err
	}

	if result == "y" {
		fileName := viper.ConfigFileUsed()
		err := os.Remove(fileName)
		if err != nil {
			fmt.Printf("Error on remove old config file %v\n", err)
			return err
		}
	}
	return nil
}
