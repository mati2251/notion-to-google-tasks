package sync

import (
	"context"
	"testing"

	"github.com/jomei/notionapi"
	"github.com/mati2251/notion-to-google-tasks/config/auth"
	"github.com/mati2251/notion-to-google-tasks/db"
	"github.com/mati2251/notion-to-google-tasks/google"
	"github.com/mati2251/notion-to-google-tasks/keys"
	"github.com/mati2251/notion-to-google-tasks/models"
	"github.com/mati2251/notion-to-google-tasks/notion"
	"github.com/mati2251/notion-to-google-tasks/test"
	"github.com/spf13/viper"
)

func TestNotionInserts(t *testing.T) {
	err := db.OpenFile()
	if err != nil {
		t.Error(err)
	}
	connection := test.GetTestConnection()
	taskDetails := test.CreateDetails()
	list, err := auth.TasksService.Tasks.List(connection.TasksListId).Do()
	ids := make([]string, 0)
	for _, task := range list.Items {
		ids = append(ids, task.Id)
	}
	if err != nil {
		t.Error(err)
	}
	taskId, _, err := google.Service.Insert(connection.NotionDatabaseId, &taskDetails)
	if err != nil {
		t.Error(err)
	}
	err = notionInserts(&ids, connection.NotionDatabaseId)
	if err != nil {
		t.Error(err)
	}
	connectedTask, err := db.GetConnectedTaskByTaskId(taskId)
	if err != nil {
		t.Error(err)
	}
	if connectedTask.TasksId != taskId {
		t.Errorf("expected %s got %s", taskId, connectedTask.TasksId)
	}
	compareTasks(t, *connectedTask)
	t.Cleanup(func() { cleanUp(t, *connectedTask) })
}

func TestGoogleInserts(t *testing.T) {
	err := db.OpenFile()

	if err != nil {
		t.Error(err)
	}
	connection := test.GetTestConnection()
	taskDetails := test.CreateDetails()
	items, err := auth.NotionClient.Database.Query(context.Background(), notionapi.DatabaseID(connection.NotionDatabaseId), &notionapi.DatabaseQueryRequest{})
	if err != nil {
		t.Error(err)
	}
	ids := make([]string, 0)
	for _, task := range items.Results {
		ids = append(ids, string(task.ID))
	}
	notionId, _, err := notion.Service.Insert(connection.NotionDatabaseId, &taskDetails)
	if err != nil {
		t.Error(err)
	}
	err = googleInserts(&ids, connection.NotionDatabaseId)
	if err != nil {
		t.Error(err)
	}
	connectedTask, err := db.GetConnectedTaskByNotionId(notionId)
	if err != nil {
		t.Error(err)
	}
	if connectedTask.NotionId != notionId {
		t.Errorf("expected %s got %s", notionId, connectedTask.NotionId)
	}
	compareTasks(t, *connectedTask)
	t.Cleanup(func() { cleanUp(t, *connectedTask) })
}

func compareTasks(t *testing.T, connectedTask models.ConnectedTask) {
	notionTask, notionUpdated, err := notion.Service.GetTaskDetails(connectedTask.ConnectionId, connectedTask.NotionId)
	if err != nil {
		t.Error(err)
	}
	googletTask, googleUpdated, err := google.Service.GetTaskDetails(connectedTask.ConnectionId, connectedTask.TasksId)
	if err != nil {
		t.Error(err)
	}
	if notionTask.Title != googletTask.Title {
		t.Errorf("expected %s got %s", notionTask.Title, googletTask.Title)
	}
	if connectedTask.TaskUpdate.Unix() != googleUpdated.Unix() {
		t.Errorf("expected %d got %d", connectedTask.TaskUpdate.Unix(), googleUpdated.Unix())
	}
	if connectedTask.NotionUpdate.Unix() != notionUpdated.Unix() {
		t.Errorf("expected %d got %d", connectedTask.TaskUpdate.Unix(), googleUpdated.Unix())
	}
}

func cleanUp(t *testing.T, connectedTask models.ConnectedTask) {
	tasksListId := viper.GetString(keys.CONNECTIONS + "." + connectedTask.ConnectionId)
	err := auth.TasksService.Tasks.Delete(tasksListId, connectedTask.TasksId).Do()
	if err != nil {
		t.Error(err)
	}
	_, err = auth.NotionClient.Page.Update(context.Background(), notionapi.PageID(connectedTask.NotionId), &notionapi.PageUpdateRequest{
		Archived:   true,
		Properties: notionapi.Properties{},
	})
	if err != nil {
		t.Error(err)
	}
	err = db.RemoveTaskByTaskId(connectedTask.TasksId)
	if err != nil {
		t.Error(err)
	}
}
