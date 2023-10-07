package google

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jomei/notionapi"
	"github.com/mati2251/notion-to-google-tasks/config/auth"
	"github.com/mati2251/notion-to-google-tasks/keys"
	"github.com/mati2251/notion-to-google-tasks/models"
	"github.com/mati2251/notion-to-google-tasks/sync/notion"
	"github.com/spf13/viper"
	"google.golang.org/api/tasks/v1"
)

func CreateNewTask(connectedTask models.ConnectedTask) string {
	title, err := notion.GetProp(*connectedTask.Notion, keys.NOTION_NAME_KEY)
	if err != nil {
		log.Fatalf("Error getting name from notion page: %v", err)
	}
	notes := createNotes(*connectedTask.Notion)
	due, _ := notion.GetProp(*connectedTask.Notion, keys.NOTION_DEADLINE_KEY)
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
	_, err = notion.UpdateValueFromProp(connectedTask.Notion, keys.TASK_ID_KEY, connectedTask.Task.Id)
	if err != nil {
		log.Fatalf("Error updating notion page: %v", err)
	}
	return task.Id
}

func createNotes(tuple notionapi.Page) string {
	notes := keys.BREAK_LINE
	notes += notion.GetPropsToString(tuple)
	notes += fmt.Sprintf("%v: %v\n", keys.NOTION_LINK_KEY, tuple.URL)
	notes += fmt.Sprintf("%v: %v\n", keys.NOTION_ID_KEY, tuple.ID)
	return notes
}

func Update(connectedTask models.ConnectedTask) error {
	title, err := notion.GetProp(*connectedTask.Notion, keys.NOTION_NAME_KEY)
	if err != nil {
		return errors.Join(err, errors.New("error getting nontio page"))
	}
	notes := strings.Split(connectedTask.Task.Notes, keys.BREAK_LINE)
	if len(notes) > 1 {
		notes = notes[:1]
	}
	notes = append(notes, createNotes(*connectedTask.Notion))
	due, err := notion.GetProp(*connectedTask.Notion, keys.NOTION_DEADLINE_KEY)
	if err != nil {
		return errors.Join(err, errors.New("error getting deadline from notion page"))
	}
	done, err := notion.GetProp(*connectedTask.Notion, keys.NOTION_STATUS_KEY)
	if err != nil {
		return errors.Join(err, errors.New("error getting status from notion page"))
	}
	if done == viper.GetString(keys.NOTION_DONE_STATUS_VALUE) {
		connectedTask.Task.Status = "completed"
	}
	newTask := connectedTask.Task
	newTask.Title = title
	newTask.Notes = strings.Join(notes, "")
	newTask.Due = due
	_, err = auth.TasksService.Tasks.Update(connectedTask.Connection.TasksList.Id, connectedTask.Task.Id, newTask).Do()
	if err != nil {
		return errors.Join(err, errors.New("error updating task"))
	}
	return nil
}
