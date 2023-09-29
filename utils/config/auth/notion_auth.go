package auth

import (
	"errors"
	"fmt"
	"log"

	"github.com/jomei/notionapi"
	"github.com/manifoldco/promptui"
	"github.com/mati2251/notion-to-google-tasks/utils/keys"
	"github.com/spf13/viper"
)

func NotionConfig() {
	configNameAndDueTimeKeys()
	var err error
	NotionClient, err = setToken()
	if err != nil {
		log.Fatalf("Unable get Notion Token: %v\n", err)
	}
}

func GetNotionToken() (*notionapi.Client, error) {
	tokenString := viper.GetString(keys.NOTION_TOKEN_KEY)
	if tokenString == "" {
		return nil, errors.New("notion token is null")
	}
	token := notionapi.Token(viper.GetString(keys.NOTION_TOKEN_KEY))
	return notionapi.NewClient(token), nil
}

func RemoveNotionConfig() {
	viper.Set(keys.NOTION_TOKEN_KEY, "")
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
	viper.Set(keys.NOTION_TOKEN_KEY, tok)
	viper.SafeWriteConfig()
	viper.WriteConfig()
	tokenString := notionapi.Token(tok)
	return notionapi.NewClient(tokenString), nil
}

func configNameAndDueTimeKeys() {
	prompt := promptui.Prompt{
		Label: "Enter name key in notion database property",
	}
	name, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	viper.Set(keys.NOTION_NAME_KEY, name)
	prompt = promptui.Prompt{
		Label: "Enter due time key in notion database property",
	}
	dueTime, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	viper.Set(keys.NOTION_DUE_TIME_KEY, dueTime)
}
