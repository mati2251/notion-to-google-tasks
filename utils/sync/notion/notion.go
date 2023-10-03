package notion

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jomei/notionapi"
	"github.com/mati2251/notion-to-google-tasks/utils/config/auth"
	"github.com/mati2251/notion-to-google-tasks/utils/keys"
	"github.com/mati2251/notion-to-google-tasks/utils/models"
	"github.com/spf13/viper"
)

func CreateDbPropTasksIdIfNotExists(databaseId notionapi.DatabaseID) {
	result, _ := auth.NotionClient.Database.Query(context.Background(), databaseId, nil)
	if result.Results[0].Properties[keys.TASK_ID_KEY] == nil {
		createDbPropTasksId(databaseId)
	}
}

func createDbPropTasksId(databaseId notionapi.DatabaseID) {
	defaultValue := &notionapi.RichTextPropertyConfig{
		Type: "rich_text",
	}
	properties := notionapi.PropertyConfigs(map[string]notionapi.PropertyConfig{keys.TASK_ID_KEY: defaultValue})
	_, err := auth.NotionClient.Database.Update(context.Background(), databaseId, &notionapi.DatabaseUpdateRequest{
		Properties: properties,
	})
	if err != nil {
		log.Fatalf("Error updating database: %v", err)
	}
}

func GetName(page notionapi.Page) (string, error) {
	nameKey := viper.GetString(keys.NOTION_NAME_KEY)
	if page.Properties[nameKey] == nil {
		return "", errors.New("Invalid notion name key: " + nameKey)
	}
	return GetStringValueFromProperty(page.Properties[nameKey]), nil
}

func GetDeadlineForTasks(page notionapi.Page) string {
	deadlineKey := viper.GetString(keys.NOTION_DEADLINE_KEY)
	deadline := ""
	if page.Properties[deadlineKey] != nil {
		deadline = GetStringValueFromProperty(page.Properties[deadlineKey])
	}
	return deadline
}

func GetPropsToString(page notionapi.Page) string {
	deadlineKey := viper.GetString(keys.NOTION_DEADLINE_KEY)
	nameKey := viper.GetString(keys.NOTION_NAME_KEY)
	propsString := ""
	for key, value := range page.Properties {
		if key != keys.TASK_ID_KEY && key != nameKey && key != deadlineKey {
			propsString += fmt.Sprintf("%v: %v\n", key, GetStringValueFromProperty(value))
		}
	}
	return propsString
}

func InsertTaskId(connectedTask models.ConnectedTask) {
	_, err := auth.NotionClient.Page.Update(context.Background(), notionapi.PageID(connectedTask.Notion.ID), &notionapi.PageUpdateRequest{
		Properties: notionapi.Properties{
			keys.TASK_ID_KEY: &notionapi.RichTextProperty{
				RichText: []notionapi.RichText{
					{
						Type: "text",
						Text: &notionapi.Text{
							Content: connectedTask.Task.Id,
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("Error updating page: %v", err)
	}
}
