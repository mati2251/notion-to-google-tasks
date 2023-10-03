package google

import (
	"fmt"
	"log"

	"github.com/jomei/notionapi"
	"github.com/mati2251/notion-to-google-tasks/utils/config/auth"
	"github.com/mati2251/notion-to-google-tasks/utils/keys"
	"github.com/mati2251/notion-to-google-tasks/utils/models"
	"github.com/mati2251/notion-to-google-tasks/utils/sync/notion"
	"google.golang.org/api/tasks/v1"
)

func updateNotes(notionPage notionapi.Page, taskId string, connection models.Connection) string {
	// auth.TasksService.Tasks.Update(connection.TasksList.Id, taskId, &tasks.Task{
	// 	Notes: createNotes(notionPage),

	// })
	return ""
}

func CreateNewTask(connectedTask models.ConnectedTask) string {
	title, err := notion.GetName(*connectedTask.Notion)
	if err != nil {
		log.Fatalf("Error getting name from notion page: %v", err)
	}
	notes := createNotes(*connectedTask.Notion)
	due := notion.GetDeadlineForTasks(*connectedTask.Notion)
	newTask := &tasks.Task{
		Title: title,
		Notes: notes,
		Due:   due,
	}
	task, err := auth.TasksService.Tasks.Insert(connectedTask.Connection.TasksList.Id, newTask).Do()
	connectedTask.Task = task
	if err != nil {
		log.Fatalf("Error creating task: %v", err)
	}
	notion.InsertTaskId(connectedTask)
	return task.Id
}

func createNotes(tuple notionapi.Page) string {
	notes := keys.BREAK_LINE
	notes += notion.GetPropsToString(tuple)
	notes += fmt.Sprintf("%v: %v\n", keys.NOTION_LINK_KEY, tuple.URL)
	notes += fmt.Sprintf("%v: %v\n", keys.NOTION_ID_KEY, tuple.ID)
	return notes
}
