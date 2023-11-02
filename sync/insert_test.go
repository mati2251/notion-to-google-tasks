package sync

import (
	"testing"

	"github.com/mati2251/notion-to-google-tasks/config/auth"
	"github.com/mati2251/notion-to-google-tasks/db"
	"github.com/mati2251/notion-to-google-tasks/google"
	"github.com/mati2251/notion-to-google-tasks/models"
	"github.com/mati2251/notion-to-google-tasks/notion"
	"github.com/mati2251/notion-to-google-tasks/test"
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
	taskId, _, err := google.Service.Insert(connection.TasksListId, &taskDetails)
	if err != nil {
		t.Error(err)
	}
	err = notionInserts(ids, connection.TasksListId)
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
