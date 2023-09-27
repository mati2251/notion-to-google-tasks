package utils

import (
	"errors"
	"fmt"
	"log"

	"github.com/jomei/notionapi"
	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
)

const NOTION_TOKEN_KEY = "notion.token"
const NOTION_NAME_KEY = "notion.name"
const NOTION_DUE_TIME_KEY = "notion.due_time"

func NotionConfig() (*notionapi.Client, error) {
	configNameAndDueTimeKeys()
	return setToken()
}

func setToken() (*notionapi.Client, error) {
	fmt.Print("Enter your Notion token (to get this add new internal integrations from https://www.notion.so/my-integrations )")
	prompt := promptui.Prompt{
		Label: "Notion token",
	}
	tok, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return nil, err
	}
	viper.Set(NOTION_TOKEN_KEY, tok)
	viper.SafeWriteConfig()
	viper.WriteConfig()
	tokenString := notionapi.Token(tok)
	return notionapi.NewClient(tokenString), nil
}

func GetNotionToken() (*notionapi.Client, error) {
	tokenString := viper.GetString(NOTION_TOKEN_KEY)
	if tokenString == "" {
		return nil, errors.New("notion token is null")
	}
	token := notionapi.Token(viper.GetString(NOTION_TOKEN_KEY))
	return notionapi.NewClient(token), nil
}

func RemoveNotionConfig() {
	viper.Set(NOTION_TOKEN_KEY, "")
}

func configNameAndDueTimeKeys() {
	prompt := promptui.Prompt{
		Label: "Enter name key in notion database property:",
	}
	name, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	viper.Set(NOTION_NAME_KEY, name)
	prompt = promptui.Prompt{
		Label: "Enter due time key in notion database property:",
	}
	dueTime, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	viper.Set(NOTION_DUE_TIME_KEY, dueTime)
}
