package notion

import (
	"context"
	"testing"
	"time"

	"github.com/jomei/notionapi"
	"github.com/mati2251/notion-to-google-tasks/config/auth"
	"github.com/mati2251/notion-to-google-tasks/keys"
	"github.com/mati2251/notion-to-google-tasks/models"
	"github.com/mati2251/notion-to-google-tasks/test"
	"github.com/spf13/viper"
)

var connection models.Connection = test.CreateMockConnection()

func TestInsert(t *testing.T) {
	details := test.CreateDetails()
	id, updated, err := Service.Insert(connection.TasksListId, &details)
	if err != nil {
		t.Error(err)
	}
	page, err := auth.NotionClient.Page.Get(context.Background(), notionapi.PageID(id))
	if err != nil {
		t.Error(err)
	}
	newTitle := GetStringValueFromProperty(page.Properties[viper.GetString(keys.NOTION_NAME_KEY)])
	if newTitle != details.Title {
		t.Errorf("Title is not correct: new value %v correct value %v", newTitle, details.Title)
	}
	newDate := GetStringValueFromProperty(page.Properties[viper.GetString(keys.NOTION_DEADLINE_KEY)])
	if newDate != details.DueDate.Format(time.RFC3339) {
		t.Errorf("Deadline is not correct: new value %v correct value %v", newDate, details.DueDate.Format(time.RFC3339))
	}
	if page.LastEditedTime != *updated {
		t.Error("Last edited time is not correct")
	}
	t.Cleanup(func() {
		print("test")
		_, err := auth.NotionClient.Page.Update(context.Background(), notionapi.PageID(id), &notionapi.PageUpdateRequest{
			Archived:   true,
			Properties: notionapi.Properties{},
		})
		if err != nil {
			t.Error(err)
		}
	})
}
