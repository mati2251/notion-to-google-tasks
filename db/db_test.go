package db

import (
	"testing"
	"time"

	"github.com/mati2251/notion-to-google-tasks/models"
)

func TestInsertAndGet(t *testing.T) {
	err := OpenFile()
	if err != nil {
		t.Error(err)
	}
	now := time.Now()
	connectedTask := models.ConnectedTask{
		TasksId:      "task-123",
		NotionId:     "notion-123",
		TaskUpdate:   &now,
		NotionUpdate: &now,
		ConnectionId: "conn-123",
	}
	err = Insert(connectedTask)
	if err != nil {
		t.Error(err)
	}
	connectedTaskFromDb, err := GetTask(connectedTask.TasksId)
	if err != nil {
		t.Error(err)
	}
	if connectedTaskFromDb.TasksId != connectedTask.TasksId {
		t.Errorf("expected %s got %s", connectedTask.TasksId, connectedTaskFromDb.TasksId)
	}
	if connectedTaskFromDb.NotionId != connectedTask.NotionId {
		t.Errorf("expected %s got %s", connectedTask.NotionId, connectedTaskFromDb.NotionId)
	}
	if connectedTaskFromDb.ConnectionId != connectedTask.ConnectionId {
		t.Errorf("expected %s got %s", connectedTask.ConnectionId, connectedTaskFromDb.ConnectionId)
	}
	if connectedTaskFromDb.TaskUpdate.Unix() != connectedTask.TaskUpdate.Unix() {
		t.Errorf("expected %d got %d", connectedTask.TaskUpdate.Unix(), connectedTaskFromDb.TaskUpdate.Unix())
	}
	if connectedTaskFromDb.NotionUpdate.Unix() != connectedTask.NotionUpdate.Unix() {
		t.Errorf("expected %d got %d", connectedTask.NotionUpdate.Unix(), connectedTaskFromDb.NotionUpdate.Unix())
	}
	t.Cleanup(func() {
		err = RemoveTask(connectedTask.TasksId)
		connectedTaskFromDb, err = GetTask(connectedTask.TasksId)
		if err == nil {
			t.Error("Removing task failed")
		}
	})
}
