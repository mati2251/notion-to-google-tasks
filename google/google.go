package google

import (
	"errors"
	"strings"
	"time"

	"github.com/mati2251/notion-to-google-tasks/config/auth"
	"github.com/mati2251/notion-to-google-tasks/db"
	"github.com/mati2251/notion-to-google-tasks/keys"
	"github.com/mati2251/notion-to-google-tasks/models"
	"github.com/spf13/viper"
	"google.golang.org/api/tasks/v1"
)

type GoogleTaskService struct{}

var Service models.Service = GoogleTaskService{}

func (GoogleTaskService) Insert(connectionId string, details *models.TaskDetails) (string, *time.Time, error) {
	done := "needsAction"
	if details.Done {
		done = "completed"
	}
	newTask := &tasks.Task{
		Title:  details.Title,
		Notes:  details.Notes,
		Status: done,
	}
	if details.DueDate != nil {
		newTask.Due = details.DueDate.Format(time.RFC3339)
	}

	taskListId := viper.GetString(keys.CONNECTIONS + "." + connectionId)
	task, err := auth.TasksService.Tasks.Insert(taskListId, newTask).Do()
	if err != nil {
		return "", nil, errors.Join(err, errors.New("error while inserting task"))
	}
	updated, err := time.Parse(time.RFC3339, task.Updated)
	if err != nil {
		return "", nil, errors.Join(err, errors.New("error while parsing updated date"))
	}
	return task.Id, &updated, nil
}

func (GoogleTaskService) Update(connectionId string, id string, details *models.TaskDetails) (*time.Time, error) {
	done := "needsAction"
	if details.Done {
		done = "completed"
	}
	dueDate := ""
	if details.DueDate != nil {
		dueDate = details.DueDate.Format(time.RFC3339)
	}
	taskListId := viper.GetString(keys.CONNECTIONS + "." + connectionId)
	task, err := auth.TasksService.Tasks.Update(taskListId, id, &tasks.Task{
		Id:     id,
		Title:  details.Title,
		Due:    dueDate,
		Notes:  details.Notes,
		Status: done,
	}).Do()
	if err != nil {
		return nil, errors.Join(err, errors.New("error while updating task"))
	}
	if details.Done {
		db.RemoveTaskByTaskId(id)
	}
	updated, err := time.Parse(time.RFC3339, task.Updated)
	if err != nil {
		return nil, errors.Join(err, errors.New("error while parsing updated date"))
	}
	return &updated, nil
}

func (GoogleTaskService) GetTaskDetails(connectionId string, id string) (*models.TaskDetails, *time.Time, error) {
	taskListId := viper.GetString(keys.CONNECTIONS + "." + connectionId)
	task, err := auth.TasksService.Tasks.Get(taskListId, id).Do()
	if err != nil {
		return nil, nil, errors.Join(err, errors.New("error while getting task"))
	}
	notesArr := strings.Split(task.Notes, keys.BREAK_LINE)
	notes := notesArr[0]
	if len(notesArr) < 1 {
		notes = ""
	}
	dueDate, err := time.Parse(time.RFC3339, task.Due)
	dueDatePointer := &dueDate
	if err != nil {
		if task.Due == "" {
			dueDatePointer = nil
		} else {
			return nil, nil, errors.Join(err, errors.New("error while parsing due date"))
		}
	}
	updated, err := time.Parse(time.RFC3339, task.Updated)
	if err != nil {
		return nil, nil, errors.Join(err, errors.New("error while updating task"))
	}
	taskDetails := &models.TaskDetails{
		Title:   task.Title,
		Done:    task.Status == "completed",
		Notes:   notes,
		DueDate: dueDatePointer,
	}
	return taskDetails, &updated, nil
}

func (GoogleTaskService) GetConnectedTaskById(id string) (*models.ConnectedTask, error) {
	return db.GetConnectedTaskByTaskId(id)
}
