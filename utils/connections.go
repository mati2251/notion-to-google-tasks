package utils

import (
	"net/http"

	"github.com/jomei/notionapi"
	"github.com/manifoldco/promptui"
)

func ConfigConnections(googleClient *http.Client, notionClient *notionapi.Client) {
	prompt := promptui.Prompt{
		Label:     "Client ID",
		AllowEdit: false,
	}

	prompt.Run()
}
