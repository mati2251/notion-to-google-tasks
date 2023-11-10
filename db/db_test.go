package db

import (
	"testing"
	"time"

	"github.com/mati2251/notion-to-google-tasks/models"
	"github.com/mati2251/notion-to-google-tasks/test"
)

func TestInsertAndGet(t *testing.T) {
	test.InitViper()
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
	connectedTaskFromDb, err := GetConnectedTaskByTaskId(connectedTask.TasksId)
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
		err = RemoveTaskByTaskId(connectedTask.TasksId)
		connectedTaskFromDb, err = GetConnectedTaskByTaskId(connectedTask.TasksId)
		if err == nil {
			t.Error("Removing task failed")
		}
	})
}

func TestGetConnectedTasks(t *testing.T) {
	test.InitViper()
	connection := test.GetTestConnection()
	err := OpenFile()
	if err != nil {
		t.Error(err)
	}

	now := time.Now()
	connectedTask1 := models.ConnectedTask{
		TasksId:      "task-123",
		NotionId:     "notion-123",
		TaskUpdate:   &now,
		NotionUpdate: &now,
		ConnectionId: connection.TasksListId,
	}
	err = Insert(connectedTask1)
	if err != nil {
		t.Error(err)
	}

	connectedTask2 := models.ConnectedTask{
		TasksId:      "task-456",
		NotionId:     "notion-456",
		TaskUpdate:   &now,
		NotionUpdate: &now,
		ConnectionId: connection.TasksListId,
	}
	err = Insert(connectedTask2)
	if err != nil {
		t.Error(err)
	}

	connectedTasks, err := GetConnectedTasks(connection.TasksListId)
	if err != nil {
		t.Error(err)
	}

	if len(connectedTasks) != 2 {
		t.Errorf("expected 2 connected tasks, got %d", len(connectedTasks))
	}

	if connectedTasks[0].TasksId != connectedTask1.TasksId {
		t.Errorf("expected %s got %s", connectedTask1.TasksId, connectedTasks[0].TasksId)
	}
	if connectedTasks[0].NotionId != connectedTask1.NotionId {
		t.Errorf("expected %s got %s", connectedTask1.NotionId, connectedTasks[0].NotionId)
	}
	if connectedTasks[0].ConnectionId != connectedTask1.ConnectionId {
		t.Errorf("expected %s got %s", connectedTask1.ConnectionId, connectedTasks[0].ConnectionId)
	}
	if connectedTasks[0].TaskUpdate.Unix() != connectedTask1.TaskUpdate.Unix() {
		t.Errorf("expected %d got %d", connectedTask1.TaskUpdate.Unix(), connectedTasks[0].TaskUpdate.Unix())
	}
	if connectedTasks[0].NotionUpdate.Unix() != connectedTask1.NotionUpdate.Unix() {
		t.Errorf("expected %d got %d", connectedTask1.NotionUpdate.Unix(), connectedTasks[0].NotionUpdate.Unix())
	}

	if connectedTasks[1].TasksId != connectedTask2.TasksId {
		t.Errorf("expected %s got %s", connectedTask2.TasksId, connectedTasks[1].TasksId)
	}
	if connectedTasks[1].NotionId != connectedTask2.NotionId {
		t.Errorf("expected %s got %s", connectedTask2.NotionId, connectedTasks[1].NotionId)
	}
	if connectedTasks[1].ConnectionId != connectedTask2.ConnectionId {
		t.Errorf("expected %s got %s", connectedTask2.ConnectionId, connectedTasks[1].ConnectionId)
	}
	if connectedTasks[1].TaskUpdate.Unix() != connectedTask2.TaskUpdate.Unix() {
		t.Errorf("expected %d got %d", connectedTask2.TaskUpdate.Unix(), connectedTasks[1].TaskUpdate.Unix())
	}
	if connectedTasks[1].NotionUpdate.Unix() != connectedTask2.NotionUpdate.Unix() {
		t.Errorf("expected %d got %d", connectedTask2.NotionUpdate.Unix(), connectedTasks[1].NotionUpdate.Unix())
	}
	t.Cleanup(func() {
		err = RemoveTaskByTaskId(connectedTask1.TasksId)
		if err != nil {
			t.Error(err)
		}
		err = RemoveTaskByTaskId(connectedTask2.TasksId)
		if err != nil {
			t.Error(err)
		}
	})

}

func TestUpdateNotionTime(t *testing.T) {
	err := OpenFile()
	if err != nil {
		t.Error(err)
	}
	now := time.Now()
	connectedTask1 := models.ConnectedTask{
		TasksId:      "task-123",
		NotionId:     "notion-123",
		TaskUpdate:   &now,
		NotionUpdate: &now,
		ConnectionId: "test",
	}
	err = Insert(connectedTask1)
	if err != nil {
		t.Error(err)
	}
	newTime := now.Add(time.Hour)
	err = UpdateNotionTime(connectedTask1.NotionId, &newTime)
	if err != nil {
		t.Error(err)
	}

	connectedTaskFromDb, err := GetConnectedTaskByNotionId(connectedTask1.NotionId)
	if err != nil {
		t.Error(err)
	}

	if connectedTaskFromDb.NotionUpdate.Unix() != now.Add(time.Hour).Unix() {
		t.Errorf("expected %d got %d", now.Add(time.Hour).Unix(), connectedTaskFromDb.NotionUpdate.Unix())
	}

	t.Cleanup(func() {
		err = RemoveTaskByTaskId(connectedTask1.TasksId)
		if err != nil {
			t.Error(err)
		}
	})
}

func TestUpdateGoogleTime(t *testing.T) {
	err := OpenFile()
	if err != nil {
		t.Error(err)
	}
	now := time.Now()
	connectedTask1 := models.ConnectedTask{
		TasksId:      "task-123",
		NotionId:     "notion-123",
		TaskUpdate:   &now,
		NotionUpdate: &now,
		ConnectionId: "test",
	}
	err = Insert(connectedTask1)
	if err != nil {
		t.Error(err)
	}
	newTime := now.Add(time.Hour)
	err = UpdateGoogleTime(connectedTask1.TasksId, &newTime)
	if err != nil {
		t.Error(err)
	}

	connectedTaskFromDb, err := GetConnectedTaskByTaskId(connectedTask1.TasksId)
	if err != nil {
		t.Error(err)
	}

	if connectedTaskFromDb.TaskUpdate.Unix() != newTime.Unix() {
		t.Errorf("expected %d got %d", newTime.Unix(), connectedTaskFromDb.NotionUpdate.Unix())
	}

	t.Cleanup(func() {
		err = RemoveTaskByTaskId(connectedTask1.TasksId)
		if err != nil {
			t.Error(err)
		}
	})
}
