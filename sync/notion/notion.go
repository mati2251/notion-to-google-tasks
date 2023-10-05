package notion

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jomei/notionapi"
	"github.com/mati2251/notion-to-google-tasks/config/auth"
	"github.com/mati2251/notion-to-google-tasks/keys"
	"github.com/mati2251/notion-to-google-tasks/models"
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

func GetProp(page notionapi.Page, viperKey string) (string, error) {
	key := viper.GetString(viperKey)
	if page.Properties[key] == nil {
		return "", errors.New("Invalid notion key: " + key)
	}
	return GetStringValueFromProperty(page.Properties[key]), nil
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

func Update(connectedTask models.ConnectedTask) error {
	newTitle := connectedTask.Task.Title
	newDue := connectedTask.Task.Due
	done := connectedTask.Task.Status == "completed"
	var err error = nil
	if done {
		doneErr := UpdateValueFromProp(connectedTask.Notion, viper.GetString(keys.NOTION_STATUS_KEY), viper.GetString(keys.NOTION_DONE_STATUS_VALUE))
		if doneErr != nil {
			err = errors.Join(err, doneErr)
		}
	}
	dueErr := UpdateValueFromProp(connectedTask.Notion, viper.GetString(keys.NOTION_DEADLINE_KEY), newDue)
	if dueErr != nil {
		err = errors.Join(err, dueErr)
	}
	titleErr := UpdateValueFromProp(connectedTask.Notion, viper.GetString(keys.NOTION_NAME_KEY), newTitle)
	if titleErr != nil {
		return errors.Join(err, titleErr)
	}
	if err != nil {
		return errors.Join(err, errors.New("error updating notion page"))
	}
	return nil
}
