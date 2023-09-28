package config

import (
	"errors"
	"fmt"
	"log"

	"github.com/jomei/notionapi"
	"github.com/manifoldco/promptui"
	"github.com/mati2251/notion-to-google-tasks/utils/basic"
	"github.com/spf13/viper"
)

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
	viper.Set(basic.NOTION_TOKEN_KEY, tok)
	viper.SafeWriteConfig()
	viper.WriteConfig()
	tokenString := notionapi.Token(tok)
	return notionapi.NewClient(tokenString), nil
}

func GetNotionToken() (*notionapi.Client, error) {
	tokenString := viper.GetString(basic.NOTION_TOKEN_KEY)
	if tokenString == "" {
		return nil, errors.New("notion token is null")
	}
	token := notionapi.Token(viper.GetString(basic.NOTION_TOKEN_KEY))
	return notionapi.NewClient(token), nil
}

func RemoveNotionConfig() {
	viper.Set(basic.NOTION_TOKEN_KEY, "")
}

func configNameAndDueTimeKeys() {
	prompt := promptui.Prompt{
		Label: "Enter name key in notion database property:",
	}
	name, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	viper.Set(basic.NOTION_NAME_KEY, name)
	prompt = promptui.Prompt{
		Label: "Enter due time key in notion database property:",
	}
	dueTime, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	viper.Set(basic.NOTION_DUE_TIME_KEY, dueTime)
}
