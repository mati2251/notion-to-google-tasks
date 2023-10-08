package google

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jomei/notionapi"
	"github.com/mati2251/notion-to-google-tasks/config/auth"
	"github.com/mati2251/notion-to-google-tasks/keys"
	"github.com/mati2251/notion-to-google-tasks/models"
	"github.com/mati2251/notion-to-google-tasks/sync/notion"
	"github.com/spf13/viper"
	"google.golang.org/api/tasks/v1"
)

func New(connectedTask models.ConnectedTask) (models.ConnectedTask, error) {
	title, err := notion.GetProp(*connectedTask.Notion, keys.NOTION_NAME_KEY)
	if err != nil {
		return connectedTask, errors.Join(err, errors.New("error getting name from notion page"))
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
		return connectedTask, errors.Join(err, errors.New("error creating task"))
	}
	page, err := notion.UpdateValueFromProp(connectedTask.Notion, keys.TASK_ID_KEY, connectedTask.Task.Id)
	if err != nil {
		return connectedTask, errors.Join(err, errors.New("error updating notion page"))
	}
	connectedTask.Task = task
	connectedTask.Notion = page
	return connectedTask, nil
}

func createNotes(tuple notionapi.Page) string {
	notes := keys.BREAK_LINE
	notes += notion.GetPropsToString(tuple)
	notes += fmt.Sprintf("%v: %v\n", keys.NOTION_LINK_KEY, tuple.URL)
	notes += fmt.Sprintf("%v: %v\n", keys.NOTION_ID_KEY, tuple.ID)
	return notes
}

func Update(connectedTask models.ConnectedTask) (models.ConnectedTask, error) {
	title, err := notion.GetProp(*connectedTask.Notion, keys.NOTION_NAME_KEY)
	if err != nil {
		return connectedTask, errors.Join(err, errors.New("error getting nontio page"))
	}
	notes := strings.Split(connectedTask.Task.Notes, keys.BREAK_LINE)
	if len(notes) > 1 {
		notes = notes[:1]
	}
	notes = append(notes, createNotes(*connectedTask.Notion))
	due, err := notion.GetProp(*connectedTask.Notion, keys.NOTION_DEADLINE_KEY)
	if err != nil {
		return connectedTask, errors.Join(err, errors.New("error getting deadline from notion page"))
	}
	done, err := notion.GetProp(*connectedTask.Notion, keys.NOTION_STATUS_KEY)
	if err != nil {
		return connectedTask, errors.Join(err, errors.New("error getting status from notion page"))
	}
	if done == viper.GetString(keys.NOTION_DONE_STATUS_VALUE) {
		connectedTask.Task.Status = "completed"
	}
	newTask := connectedTask.Task
	newTask.Title = title
	newTask.Notes = strings.Join(notes, "")
	newTask.Due = due
	task, err := auth.TasksService.Tasks.Update(connectedTask.Connection.TasksList.Id, connectedTask.Task.Id, newTask).Do()
	connectedTask.Task = task
	if err != nil {
		return connectedTask, errors.Join(err, errors.New("error updating task"))
	}
	return connectedTask, nil
}

func UpdateNotes(connectedTask models.ConnectedTask) (models.ConnectedTask, error) {
	notes := createNotes(*connectedTask.Notion)
	newTask := connectedTask.Task
	newTask.Notes = notes
	task, err := auth.TasksService.Tasks.Update(connectedTask.Connection.TasksList.Id, connectedTask.Task.Id, newTask).Do()
	connectedTask.Task = task
	if err != nil {
		return connectedTask, errors.Join(err, errors.New("error updating task"))
	}
	return connectedTask, nil
}

func SetDone(tasksList tasks.TaskList, task *tasks.Task) (*tasks.Task, error) {
	task.Status = "completed"
	return auth.TasksService.Tasks.Update(tasksList.Id, task.Id, task).Do()
}
