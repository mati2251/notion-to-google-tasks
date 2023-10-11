package cmd

import (
	"fmt"
	"os"

	"github.com/mati2251/notion-to-google-tasks/keys"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "notion-to-google-tasks",
	Short: "Synchronization Notion page to Google Tasks",
	Long: `Sync your selected sites from Notion to Google Tasks.
Config your sync by subcommand config.
Add entry to crontab to run this command periodically.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/notion-to-google-tasks/config.yaml)")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		configPath := home + keys.FILES_PATH
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			os.MkdirAll(configPath, 0755)
		}
		viper.AddConfigPath(configPath)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
