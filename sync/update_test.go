package sync

import (
	"context"
	"testing"
	"time"

	"github.com/jomei/notionapi"
	"github.com/mati2251/notion-to-google-tasks/config/auth"
	"github.com/mati2251/notion-to-google-tasks/keys"
	"github.com/mati2251/notion-to-google-tasks/models"
	"github.com/mati2251/notion-to-google-tasks/sync/google"
	"github.com/mati2251/notion-to-google-tasks/sync/notion"
	"github.com/mati2251/notion-to-google-tasks/test"
	"github.com/spf13/viper"
	"google.golang.org/api/tasks/v1"
)

var connection = test.CreateMockConnection()

func TestUpdateForceNotion(t *testing.T) {
	connectedTask, err := mockConnectedTask()
	if err != nil {
		t.Error(err)
	}
	connectedTask.Task.Title = "testTitle2"
	connectedTask.Task.Notes = "testNotes2"
	connectedTask.Task.Due = "2030-01-02T00:00:00Z"
	connectedTask, err = update(connectedTask, true)
	if err != nil {
		t.Error(err)
	}
	name, err := notion.GetProp(*connectedTask.Notion, keys.NOTION_NAME_KEY)
	if err != nil {
		t.Error(err)
	}
	if name != connectedTask.Task.Title {
		t.Error("notion name was not updated")
	}
	due, err := notion.GetProp(*connectedTask.Notion, keys.NOTION_DEADLINE_KEY)
	if err != nil {
		t.Error(err)
	}
	if due != connectedTask.Task.Due {
		t.Error("notion due was not updated")
	}
	connectedTask.Task.Status = "completed"
	connectedTask, err = update(connectedTask, true)
	if err != nil {
		t.Error(err)
	}
	status, err := notion.GetProp(*connectedTask.Notion, keys.NOTION_STATUS_KEY)
	if err != nil {
		t.Error(err)
	}
	if status != viper.GetString(keys.NOTION_DONE_STATUS_VALUE) {
		t.Error("notion status was not updated")
	}
	connectedTask.Task.Status = "needsAction"
	connectedTask, err = update(connectedTask, true)
	if err != nil {
		t.Error(err)
	}
	status, err = notion.GetProp(*connectedTask.Notion, keys.NOTION_STATUS_KEY)
	if err != nil {
		t.Error(err)
	}
	if status != viper.GetString(keys.NOTION_DONE_STATUS_VALUE) {
		t.Error("notion status was not updated")
	}
	t.Cleanup(func() {
		cleanUpMockConnectedTask(connectedTask)
	})
}

func TestUpdateForceNotionDeleteGoogle(t *testing.T) {
	connectedTask, err := mockConnectedTask()
	if err != nil {
		t.Error(err)
	}
	auth.TasksService.Tasks.Delete(connectedTask.Connection.TasksList.Id, connectedTask.Task.Id).Do()
	connectedTask, err = update(connectedTask, true)
	if err != nil {
		t.Error(err)
	}
	status, err := notion.GetProp(*connectedTask.Notion, keys.NOTION_STATUS_KEY)
	if err != nil {
		t.Error(err)
	}
	if status != viper.GetString(keys.NOTION_DONE_STATUS_VALUE) {
		t.Error("notion status was not updated")
	}
	t.Cleanup(func() {
		cleanUpMockConnectedTask(connectedTask)
	})
}

func TestUpdateForceGoogle(t *testing.T) {
	connectedTask, err := mockConnectedTask()
	if err != nil {
		t.Error(err)
	}
	page, err := notion.UpdateValueFromProp(connectedTask.Notion, viper.GetString(keys.NOTION_NAME_KEY), "testTitle2")
	connectedTask.Notion = page
	if err != nil {
		t.Error(err)
	}
	page, err = notion.UpdateValueFromProp(connectedTask.Notion, viper.GetString(keys.NOTION_DEADLINE_KEY), "2030-01-02T00:00:00Z")
	connectedTask.Notion = page
	if err != nil {
		t.Error(err)
	}
	page, err = notion.UpdateValueFromProp(connectedTask.Notion, viper.GetString(keys.NOTION_STATUS_KEY), viper.GetString(keys.NOTION_DONE_STATUS_VALUE))
	connectedTask.Notion = page
	if err != nil {
		t.Error(err)
	}
	connectedTask.Notion = page
	time.Sleep(1 * time.Second)
	connectedTask, err = update(connectedTask, true)
	if err != nil {
		t.Error(err)
	}
	if connectedTask.Task.Title != "testTitle2" {
		t.Error("google title was not updated")
	}
	if connectedTask.Task.Due != "2030-01-02T00:00:00.000Z" {
		t.Error("google due was not updated")
	}
	if connectedTask.Task.Status != "completed" {
		t.Error("google status was not updated")
	}
	t.Cleanup(func() {
		cleanUpMockConnectedTask(connectedTask)
	})
}

func TestUpdateForceGoogleDeleteNotion(t *testing.T) {
	connectedTask, err := mockConnectedTask()
	if err != nil {
		t.Error(err)
	}
	page, err := auth.NotionClient.Page.Update(context.Background(), notionapi.PageID(connectedTask.Notion.ID), &notionapi.PageUpdateRequest{
		Archived: true,
	})
	connectedTask.Notion = page
	connectedTask, err = update(connectedTask, true)
	if err != nil {
		t.Error(err)
	}
	if connectedTask.Task.Status != "completed" {
		t.Error("google status was not updated")
	}
	t.Cleanup(func() {
		cleanUpMockConnectedTask(connectedTask)
	})
}

func mockConnectedTask() (models.ConnectedTask, error) {
	new := &tasks.Task{
		Title: "testTitle",
		Id:    "testId",
		Due:   "2030-01-01T00:00:00Z",
		Notes: "testNotes",
	}
	connectedTask := models.ConnectedTask{
		Task:       new,
		Notion:     nil,
		Connection: &connection,
	}
	new.Due = "2021-01-01T00:00:00.000Z"
	connectedTask, err := notion.New(connectedTask)
	if err != nil {
		return connectedTask, err
	}
	connectedTask, err = google.New(connectedTask)
	if err != nil {
		return connectedTask, err
	}
	return connectedTask, nil
}

func cleanUpMockConnectedTask(connectedTask models.ConnectedTask) {
	if connectedTask.Task != nil {
		auth.TasksService.Tasks.Delete(connectedTask.Connection.TasksList.Id, connectedTask.Task.Id).Do()
	}
	if connectedTask.Notion != nil {
		auth.NotionClient.Page.Update(context.Background(), notionapi.PageID(connectedTask.Notion.ID), &notionapi.PageUpdateRequest{
			Archived: true,
		})
	}
}
