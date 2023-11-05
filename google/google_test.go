package google

import (
	"testing"

	"github.com/mati2251/notion-to-google-tasks/config/auth"
	"github.com/mati2251/notion-to-google-tasks/models"
	"github.com/mati2251/notion-to-google-tasks/test"
	"google.golang.org/api/tasks/v1"
)

func TestInsert(t *testing.T) {
	connection := test.GetTestConnection()
	details := test.CreateDetails()
	taskId, _, err := Service.Insert(connection.NotionDatabasId, &details)
	if err != nil {
		t.Error(err)
	}
	task, err := auth.TasksService.Tasks.Get(connection.TasksListId, taskId).Do()
	if err != nil {
		t.Error(err)
	}
	assertTask(t, *task, details)
	t.Cleanup(func() { cleanUp(t, connection, task) })
}

func TestInsertWithNullDate(t *testing.T) {
	connection := test.GetTestConnection()
	details := test.CreateDetails()
	details.DueDate = nil
	taskId, _, err := Service.Insert(connection.NotionDatabasId, &details)
	if err != nil {
		t.Error(err)
	}
	task, err := auth.TasksService.Tasks.Get(connection.TasksListId, taskId).Do()
	if err != nil {
		t.Error(err)
	}
	assertTask(t, *task, details)
	t.Cleanup(func() { cleanUp(t, connection, task) })
}

func TestUpdate(t *testing.T) {
	connection := test.GetTestConnection()
	details := test.CreateDetails()
	taskId, _, err := Service.Insert(connection.NotionDatabasId, &details)
	if err != nil {
		t.Error(err)
	}
	details.Title = "Updated title"
	newDate := details.DueDate.AddDate(0, 0, 1)
	details.DueDate = &newDate
	details.Notes = "Updated notes"
	_, err = Service.Update(connection.NotionDatabasId, taskId, &details)
	if err != nil {
		t.Error(err)
	}
	task, err := auth.TasksService.Tasks.Get(connection.TasksListId, taskId).Do()
	if err != nil {
		t.Error(err)
	}
	assertTask(t, *task, details)
	t.Cleanup(func() { cleanUp(t, connection, task) })
}

func TestGetTaskDetails(t *testing.T) {
	connection := test.GetTestConnection()
	details := test.CreateDetails()
	taskId, _, err := Service.Insert(connection.NotionDatabasId, &details)
	if err != nil {
		t.Error(err)
	}
	taskDetails, _, err := Service.GetTaskDetails(connection.NotionDatabasId, taskId)
	if err != nil {
		t.Error(err)
	}
	task, err := auth.TasksService.Tasks.Get(connection.TasksListId, taskId).Do()
	if err != nil {
		t.Error(err)
	}
	assertTask(t, *task, *taskDetails)
	t.Cleanup(func() { cleanUp(t, connection, &tasks.Task{Id: taskId}) })
}

func TestGetTaskDetailsWithNullDate(t *testing.T) {
	connection := test.GetTestConnection()
	details := test.CreateDetails()
	details.DueDate = nil
	taskId, _, err := Service.Insert(connection.NotionDatabasId, &details)
	if err != nil {
		t.Error(err)
	}
	taskDetails, _, err := Service.GetTaskDetails(connection.NotionDatabasId, taskId)
	if err != nil {
		t.Error(err)
	}
	task, err := auth.TasksService.Tasks.Get(connection.TasksListId, taskId).Do()
	if err != nil {
		t.Error(err)
	}
	assertTask(t, *task, *taskDetails)
	t.Cleanup(func() { cleanUp(t, connection, &tasks.Task{Id: taskId}) })
}

func cleanUp(t *testing.T, connection models.Connection, task *tasks.Task) {
	err := auth.TasksService.Tasks.Delete(connection.TasksListId, task.Id).Do()
	if err != nil {
		t.Error(err)
	}
}

func assertTask(t *testing.T, task tasks.Task, details models.TaskDetails) {
	if task.Title != details.Title {
		t.Errorf("Expected title: %s, got: %s", details.Title, task.Title)
	}
	if task.Notes != details.Notes {
		t.Errorf("Expected notes: %s, got: %s", details.Notes, task.Notes)
	}
	if task.Status != "needsAction" {
		t.Errorf("Expected status: needsAction, got: %s", task.Status)
	}
	if task.Due != "" && details.DueDate != nil {
		if task.Due[:10] != details.DueDate.Format("2006-01-02") {
			t.Errorf("Expected due date: %s, got: %s", details.DueDate.Format("2006-01-02"), task.Due[:10])
		}
	}
}
