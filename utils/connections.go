package utils

import (
	"net/http"

	"github.com/jomei/notionapi"
	"github.com/manifoldco/promptui"
)

func ConfigConnections(googleClient *http.Client, notionClient *notionapi.Client) {
	prompt := promptui.Prompt{
		Label:     "Share notion pages which you want synchronize",
		AllowEdit: false,
	}
	prompt.Run()
}
