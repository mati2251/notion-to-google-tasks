package utils

import (
	"fmt"

	"github.com/jomei/notionapi"
	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
)

const NOTION_TOKEN_KEY = "notion.token"

func NotionConfig() {
	setToken()
}

func setToken() {
	fmt.Print("Enter your Notion token (to get this add new internal integrations from https://www.notion.so/my-integrations )")
	prompt := promptui.Prompt{
		Label: "Notion token",
	}
	tok, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	viper.Set(NOTION_TOKEN_KEY, tok)
}

func GetNotionToken() *notionapi.Client {
	token := notionapi.Token(viper.GetString(NOTION_TOKEN_KEY))
	return notionapi.NewClient(token)
}
