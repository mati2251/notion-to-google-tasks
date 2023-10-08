package notion

import (
	"context"
	"testing"

	"github.com/jomei/notionapi"
	"github.com/mati2251/notion-to-google-tasks/config/auth"
	"github.com/mati2251/notion-to-google-tasks/keys"
	"github.com/mati2251/notion-to-google-tasks/models"
	"github.com/mati2251/notion-to-google-tasks/test"
	"github.com/spf13/viper"
	"google.golang.org/api/tasks/v1"
)

var connection models.Connection = test.CreateMockConnection()

func TestCreateAndRemoveProp(t *testing.T) {
	err := CreateProp(*connection.NotionDatabase, "test", "rich_text")
	if err != nil {
		t.Error(err)
	}
	database, err := auth.NotionClient.Database.Get(context.Background(), notionapi.DatabaseID(connection.NotionDatabase.ID))
	if err != nil {
		t.Error(err)
	}
	if database.Properties["test"] == nil {
		t.Error("property not created")
	}
	err = RemoveProp(*database, "test")
	if err != nil {
		t.Error(err)
	}
	database, err = auth.NotionClient.Database.Get(context.Background(), notionapi.DatabaseID(connection.NotionDatabase.ID))
	if err != nil {
		t.Error(err)
	}
	if database.Properties["test"] != nil {
		t.Error("property not removed")
	}
	t.Cleanup(func() {
		database, _ = auth.NotionClient.Database.Get(context.Background(), notionapi.DatabaseID(connection.NotionDatabase.ID))
		RemoveProp(*database, "test")
	})
}

func TestNewGetProp(t *testing.T) {
	conn := models.ConnectedTask{
		Task: &tasks.Task{
			Id:    "testId",
			Title: "testTitle",
			Due:   "2021-01-01T00:00:00Z",
		},
		Notion:     nil,
		Connection: &connection,
	}
	conn, err := New(conn)
	page := conn.Notion
	if err != nil {
		t.Error(err)
	}
	if GetStringValueFromProperty(page.Properties[keys.TASK_ID_KEY]) != conn.Task.Id {
		t.Error("Bad notion new tuple task id")
	}
	if GetStringValueFromProperty(page.Properties[viper.GetString(keys.NOTION_NAME_KEY)]) != conn.Task.Title {
		t.Error("Bad notion new tuple title")
	}
	test := GetStringValueFromProperty(page.Properties[viper.GetString(keys.NOTION_DEADLINE_KEY)])
	if test != conn.Task.Due {
		t.Error("Bad notion new tuple deadline")
	}
	conn.Task.Title = "newTitle"
	conn.Task.Due = "2021-01-02T00:00:00Z"
	conn.Task.Status = "completed"
	conn, err = Update(conn)
	page = conn.Notion
	if err != nil {
		t.Error(err)
	}
	if GetStringValueFromProperty(page.Properties[keys.TASK_ID_KEY]) != conn.Task.Id {
		t.Error("Bad notion update tuple task id")
	}
	if GetStringValueFromProperty(page.Properties[viper.GetString(keys.NOTION_NAME_KEY)]) != conn.Task.Title {
		t.Error("Bad notion update tuple title")
	}
	test = GetStringValueFromProperty(page.Properties[viper.GetString(keys.NOTION_DEADLINE_KEY)])
	if test != conn.Task.Due {
		t.Error("Bad notion update tuple deadline")
	}
	t.Cleanup(func() {
		_, err := auth.NotionClient.Page.Update(context.Background(), notionapi.PageID(page.ID), &notionapi.PageUpdateRequest{
			Archived:   true,
			Properties: notionapi.Properties{},
		})
		if err != nil {
			t.Error(err)
		}
	})
}
