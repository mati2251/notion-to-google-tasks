package sync

import (
	"testing"
	"time"

	"github.com/mati2251/notion-to-google-tasks/db"
	"github.com/mati2251/notion-to-google-tasks/google"
	"github.com/mati2251/notion-to-google-tasks/models"
	"github.com/mati2251/notion-to-google-tasks/notion"
	"github.com/mati2251/notion-to-google-tasks/test"
)

func TestSyncNotionToGoogle(t *testing.T) {
	newTask := test.CreateDetails()
	updatedTask := test.CreateDetails()
	updatedTask.Title = "Updated title"
	newDate := updatedTask.DueDate.AddDate(1, 0, 0)
	updatedTask.DueDate = &newDate
	testSync(newTask, updatedTask, false, t)
}

func TestSyncGoogleToNotion(t *testing.T) {
	newTask := test.CreateDetails()
	updatedTask := test.CreateDetails()
	updatedTask.Title = "Updated title"
	newDate := updatedTask.DueDate.AddDate(1, 0, 0)
	updatedTask.DueDate = &newDate
	testSync(newTask, updatedTask, true, t)
}

func TestSyncNotionToGoogleWithNullDate(t *testing.T) {
	newTask := test.CreateDetails()
	newTask.DueDate = nil
	updatedTask := test.CreateDetails()
	updatedTask.Title = "Updated title"
	testSync(newTask, updatedTask, false, t)
}

func TestSyncGoogleToNotionWithNullDate(t *testing.T) {
	newTask := test.CreateDetails()
	newTask.DueDate = nil
	updatedTask := test.CreateDetails()
	updatedTask.Title = "Updated title"
	testSync(newTask, updatedTask, true, t)
}

func TestSyncNotionToGoogleDone(t *testing.T) {
	newTask := test.CreateDetails()
	updatedTask := test.CreateDetails()
	updatedTask.Done = true
	updatedTask.Title = "Updated title"
	testSync(newTask, updatedTask, false, t)
}

func TestSyncGoogleToNotionDone(t *testing.T) {
	newTask := test.CreateDetails()
	updatedTask := test.CreateDetails()
	updatedTask.Done = false
	updatedTask.Title = "Updated title"
	testSync(newTask, updatedTask, true, t)
}

func testSync(newTask models.TaskDetails, updatedTask models.TaskDetails, isGoogleTrigger bool, t *testing.T) {
	test.InitViper()
	triggerService := notion.Service
	syncService := google.Service
	if isGoogleTrigger {
		triggerService = google.Service
		syncService = notion.Service
	}
	connection := test.GetTestConnection()
	id, _, err := triggerService.Insert(connection.NotionDatabaseId, &newTask)
	if err != nil {
		t.Error(err)
	}
	err = Sync([]models.Connection{connection})
	if err != nil {
		t.Error(err)
	}
	connectedTasks, err := triggerService.GetConnectedTaskById(id)
	if err != nil {
		t.Error(err)
	}
	newTaskFromService, _, err := syncService.GetTaskDetails(connection.NotionDatabaseId, getOpositeTaskId(isGoogleTrigger, connectedTasks))
	if err != nil {
		t.Error(err)
	}
	time.Sleep(1 * time.Minute)
	compareTaskDetails(newTask, *newTaskFromService, t)
	copyDb := &db.DB
	db.DB = nil
	_, err = triggerService.Update(connection.NotionDatabaseId, id, &updatedTask)
	db.DB = *copyDb
	if err != nil {
		t.Error(err)
	}
	err = Sync([]models.Connection{connection})
	if err != nil {
		t.Error(err)
	}
	updatedTaskFromService, _, err := syncService.GetTaskDetails(connection.NotionDatabaseId, getOpositeTaskId(isGoogleTrigger, connectedTasks))
	if err != nil {
		t.Error(err)
	}
	compareTaskDetails(updatedTask, *updatedTaskFromService, t)
}

func compareTaskDetails(task1 models.TaskDetails, task2 models.TaskDetails, t *testing.T) {
	if task1.Title != task2.Title {
		t.Errorf("Title is not equal: %s != %s", task1.Title, task2.Title)
	}
	if task1.Done != task2.Done {
		t.Errorf("Done is not equal %t != %t", task1.Done, task2.Done)
	}
	if task1.DueDate == nil && task2.DueDate == nil {
		return
	}
	if task1.DueDate.Format(time.DateOnly) != task2.DueDate.Format(time.DateOnly) {
		t.Errorf("DueDate is not equal %s != %s", task1.DueDate.Format(time.DateOnly), task2.DueDate.Format(time.DateOnly))
	}
}

func getOpositeTaskId(isGoogleTask bool, connected *models.ConnectedTask) string {
	if isGoogleTask {
		return connected.NotionId
	}
	return connected.TasksId
}
