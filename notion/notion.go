package notion

import (
	"context"
	"errors"
	"slices"

	"github.com/jomei/notionapi"
	"github.com/mati2251/notion-to-google-tasks/config/auth"
	"github.com/mati2251/notion-to-google-tasks/keys"
	"github.com/spf13/viper"
)

type NotionService struct{}

var Service NotionService

func (_ NotionService) Inserts(ids []string, connectionId string) error {
	notionId := notionapi.DatabaseID(viper.GetString(keys.CONNECTIONS))
	items, err := auth.NotionClient.Database.Query(context.Background(), notionId, &notionapi.DatabaseQueryRequest{})
	if err != nil {
		return errors.Join(err, errors.New("error while getting database"))
	}
	for _, page := range items.Results {
		if !slices.Contains(ids, page.ID.String()) {

		}
	}
	return nil
}

func (_ NotionService) Insert(taskId string, connectionId string) error {
	return nil
}
