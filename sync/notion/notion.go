package notion

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jomei/notionapi"
	"github.com/mati2251/notion-to-google-tasks/config/auth"
	"github.com/mati2251/notion-to-google-tasks/keys"
	"github.com/mati2251/notion-to-google-tasks/models"
	"github.com/spf13/viper"
)

func CreateProp(database notionapi.Database, key string, propertyType string) error {
	if database.Properties[key] != nil {
		return errors.New("property already exists")
	}
	propertyType_ := notionapi.PropertyConfigType(propertyType)
	defaultValue := &notionapi.RichTextPropertyConfig{
		Type: propertyType_,
	}
	properties := notionapi.PropertyConfigs(map[string]notionapi.PropertyConfig{key: defaultValue})
	_, err := auth.NotionClient.Database.Update(context.Background(), notionapi.DatabaseID(database.ID), &notionapi.DatabaseUpdateRequest{
		Properties: properties,
	})
	if err != nil {
		return errors.Join(err, errors.New("error creating notion property"))
	}
	return nil
}

func RemoveProp(database notionapi.Database, key string) error {
	if database.Properties[key] == nil {
		return errors.New("property doesn't exist")
	}
	_, err := auth.NotionClient.Database.Update(context.Background(), notionapi.DatabaseID(database.ID), &notionapi.DatabaseUpdateRequest{
		Properties: notionapi.PropertyConfigs(map[string]notionapi.PropertyConfig{key: nil}),
	})
	if err != nil {
		return errors.Join(err, errors.New("error removing notion property"))
	}
	return nil
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

func Update(connectedTask models.ConnectedTask) (models.ConnectedTask, error) {
	newTitle := connectedTask.Task.Title
	newDue := connectedTask.Task.Due
	done := connectedTask.Task.Status == "completed"
	var err error = nil
	if done {
		_, doneErr := UpdateValueFromProp(connectedTask.Notion, viper.GetString(keys.NOTION_STATUS_KEY), viper.GetString(keys.NOTION_DONE_STATUS_VALUE))
		if doneErr != nil {
			err = errors.Join(err, doneErr)
		}
	}
	_, dueErr := UpdateValueFromProp(connectedTask.Notion, viper.GetString(keys.NOTION_DEADLINE_KEY), newDue)
	if dueErr != nil {
		err = errors.Join(err, dueErr)
	}
	page, titleErr := UpdateValueFromProp(connectedTask.Notion, viper.GetString(keys.NOTION_NAME_KEY), newTitle)
	if titleErr != nil {
		return connectedTask, errors.Join(err, titleErr)
	}
	if err != nil {
		return connectedTask, errors.Join(err, errors.New("error updating notion page"))
	}
	connectedTask.Notion = page
	return connectedTask, nil
}

func New(connectedTask models.ConnectedTask) (models.ConnectedTask, error) {
	newTitle := connectedTask.Task.Title
	newDue, _ := time.Parse(time.RFC3339, connectedTask.Task.Due)
	newDueNotion := notionapi.Date(newDue)
	newTasksId := connectedTask.Task.Id
	Properties := notionapi.Properties{
		viper.GetString(keys.NOTION_NAME_KEY): notionapi.TitleProperty{
			Type:  "title",
			Title: []notionapi.RichText{{Type: "text", Text: &notionapi.Text{Content: newTitle}}},
		},
		viper.GetString(keys.NOTION_DEADLINE_KEY): notionapi.DateProperty{
			Type: "date",
			Date: &notionapi.DateObject{
				Start: &newDueNotion,
			},
		},
		keys.TASK_ID_KEY: NewRichTextProperty(newTasksId),
	}
	page, err := auth.NotionClient.Page.Create(context.Background(), &notionapi.PageCreateRequest{
		Properties: Properties,
		Parent: notionapi.Parent{
			Type:       "database_id",
			DatabaseID: notionapi.DatabaseID(connectedTask.Connection.NotionDatabase.ID),
		},
	})
	if err != nil {
		return connectedTask, errors.Join(err, errors.New("error creating notion tuple"))
	}
	connectedTask.Notion = page
	return connectedTask, nil
}

func SetDone(page *notionapi.Page) (*notionapi.Page, error) {
	return UpdateValueFromProp(page, viper.GetString(keys.NOTION_STATUS_KEY), viper.GetString(keys.NOTION_DONE_STATUS_VALUE))
}
