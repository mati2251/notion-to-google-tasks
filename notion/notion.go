package notion

import (
	"context"
	"errors"
	"time"

	"github.com/jomei/notionapi"
	"github.com/mati2251/notion-to-google-tasks/config/auth"
	"github.com/mati2251/notion-to-google-tasks/db"
	"github.com/mati2251/notion-to-google-tasks/keys"
	"github.com/mati2251/notion-to-google-tasks/models"
	"github.com/spf13/viper"
)

type NotionService struct{}

var Service models.Service = NotionService{}

func (NotionService) Insert(connectionId string, details *models.TaskDetails) (string, *time.Time, error) {
	databaseId := notionapi.DatabaseID(connectionId)
	properties := notionapi.Properties{
		viper.GetString(keys.NOTION_NAME_KEY): notionapi.TitleProperty{
			Title: []notionapi.RichText{NewRichText(details.Title)},
		},
	}
	if details.DueDate != nil {
		properties[viper.GetString(keys.NOTION_DEADLINE_KEY)] = NewDateProperty(*details.DueDate)
	}
	page, err := auth.NotionClient.Page.Create(context.Background(), &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			DatabaseID: databaseId,
			Type:       "database_id",
		},
		Properties: properties,
	})
	if err != nil {
		return "", nil, err
	}
	return page.ID.String(), &page.LastEditedTime, nil
}

func (NotionService) GetTaskDetails(connectionId string, id string) (*models.TaskDetails, *time.Time, error) {
	page, err := auth.NotionClient.Page.Get(context.Background(), notionapi.PageID(id))
	if err != nil {
		return nil, nil, errors.Join(err, errors.New("error while getting page"))
	}
	title := GetStringValueFromProperty(page.Properties[viper.GetString(keys.NOTION_NAME_KEY)])
	dueDateStr := GetStringValueFromProperty(page.Properties[viper.GetString(keys.NOTION_DEADLINE_KEY)])
	dueDate, err := time.Parse(time.RFC3339, dueDateStr)
	dueDatePtr := &dueDate
	done := GetStringValueFromProperty(page.Properties[viper.GetString(keys.NOTION_STATUS_KEY)]) == viper.GetString(keys.NOTION_DONE_STATUS_VALUE)
	if page.Archived {
		done = true
	}
	if err != nil {
		if dueDateStr == "" {
			dueDatePtr = nil
		} else {
			return nil, nil, errors.Join(err, errors.New("error parsing due date"))
		}
	}
	notes := createNotes(page.Properties)
	return &models.TaskDetails{
		Title:   title,
		DueDate: dueDatePtr,
		Done:    done,
		Notes:   notes,
	}, &page.LastEditedTime, nil
}

func createNotes(properties notionapi.Properties) string {
	notes := keys.BREAK_LINE
	for key, value := range properties {
		if key != viper.GetString(keys.NOTION_NAME_KEY) && key != viper.GetString(keys.NOTION_DEADLINE_KEY) {
			notes += key + ": " + GetStringValueFromProperty(value) + "\n"
		}
	}
	return notes
}

func (NotionService) Update(connectionId string, id string, details *models.TaskDetails) (*time.Time, error) {
	page, err := auth.NotionClient.Page.Get(context.Background(), notionapi.PageID(id))
	if err != nil {
		return nil, err
	}
	titleProperty, err := UpdateValueFromProp(page.Properties[viper.GetString(keys.NOTION_NAME_KEY)], details.Title)
	if err != nil {
		return nil, err
	}
	dateProperty := page.Properties[viper.GetString(keys.NOTION_DEADLINE_KEY)]
	if details.DueDate != nil {
		dateProperty, err = UpdateValueFromProp(dateProperty, details.DueDate.Format(time.RFC3339))
		if err != nil {
			return nil, err
		}
	}
	if details.Done {
		db.RemoveTaskByNotionId(id)
	}
	page, err = auth.NotionClient.Page.Update(context.Background(), notionapi.PageID(id), &notionapi.PageUpdateRequest{
		Properties: notionapi.Properties{
			viper.GetString(keys.NOTION_NAME_KEY):     titleProperty,
			viper.GetString(keys.NOTION_DEADLINE_KEY): dateProperty,
		},
	})
	if err != nil {
		return nil, err
	}
	if details.Done {
		prop := page.Properties[viper.GetString(keys.NOTION_STATUS_KEY)]
		newProp, err := UpdateValueFromProp(prop, viper.GetString(keys.NOTION_DONE_STATUS_VALUE))
		if err != nil {
			return nil, err
		}
		page, err = auth.NotionClient.Page.Update(context.Background(), notionapi.PageID(id), &notionapi.PageUpdateRequest{
			Properties: notionapi.Properties{
				viper.GetString(keys.NOTION_STATUS_KEY): newProp,
			},
		})
		return &page.LastEditedTime, err
	}
	return &page.LastEditedTime, nil
}

func (NotionService) GetConnectedTaskById(id string) (*models.ConnectedTask, error) {
	return db.GetConnectedTaskByNotionId(id)
}
