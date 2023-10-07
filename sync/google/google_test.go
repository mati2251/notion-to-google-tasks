package google

import (
	"testing"

	"github.com/mati2251/notion-to-google-tasks/config/auth"
	"github.com/mati2251/notion-to-google-tasks/models"
	"github.com/mati2251/notion-to-google-tasks/sync/notion"
	"github.com/mati2251/notion-to-google-tasks/test"
	"google.golang.org/api/tasks/v1"
)

var connection models.Connection = test.CreateMockConnection()

func TestCreateNewTask(t *testing.T) {
	conn, task := createMockPage()
	newTask := CreateNewTask(conn)
	if newTask.Title != task.Title {
		t.Errorf("Expected title to be %v, got %v", task.Title, newTask.Title)
	}
	if newTask.Due != task.Due {
		t.Errorf("Expected due to be %v, got %v", task.Due, newTask.Due)
	}
	t.Cleanup(func() {
		err := auth.TasksService.Tasks.Delete(connection.TasksList.Id, newTask.Id).Do()
		if err != nil {
			t.Error(err)
		}
	})
}

func createMockPage() (models.ConnectedTask, *tasks.Task) {
	new := &tasks.Task{
		Title: "testTitle",
		Id:    "testId",
		Due:   "2021-01-01T00:00:00Z",
		Notes: "testNotes",
	}
	connectedTask := models.ConnectedTask{
		Task:       new,
		Notion:     nil,
		Connection: &connection,
	}
	new.Due = "2021-01-01T00:00:00.000Z"
	page, err := notion.New(connectedTask)
	if err != nil {
		panic(err)
	}
	connectedTask.Notion = page
	connectedTask.Task = nil
	return connectedTask, new
}
