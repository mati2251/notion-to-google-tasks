package utils

import (
	"fmt"

	"github.com/spf13/viper"
)

const NOTION_TOKEN_KEY = "notion.token"

func NotionConfig() {
	setToken()
}

func setToken() {
	fmt.Print("Enter your Notion token (to get this add new internal integrations from https://www.notion.so/my-integrations ): ")
	var token string
	fmt.Scanln(&token)
	viper.Set(NOTION_TOKEN_KEY, token)
}

func GetNotionToken() string {
	return viper.GetString(NOTION_TOKEN_KEY)
}
