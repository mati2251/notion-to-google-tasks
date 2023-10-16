package google

import (
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/mati2251/notion-to-google-tasks/config/auth"
	"github.com/mati2251/notion-to-google-tasks/db"
	"github.com/mati2251/notion-to-google-tasks/keys"
	"github.com/mati2251/notion-to-google-tasks/models"
	"github.com/mati2251/notion-to-google-tasks/notion"
	"google.golang.org/api/tasks/v1"
)

type GoogleTaskService struct{}

var Service models.Service = GoogleTaskService{}

func (_ GoogleTaskService) Inserts(ids []string, connectionId string) error {
	list, err := auth.TasksService.Tasks.List(connectionId).Do()
	if err != nil {
		return errors.Join(err, errors.New("error while getting tasklist"))
	}
	for _, task := range list.Items {
		if !slices.Contains(ids, task.Id) {
			taskDetails, err := Service.GetTaskDetails(connectionId, task.Id)
			if err != nil {
				return err
			}
			notionId, updated, err := notion.Service.Insert(connectionId, taskDetails)
			if err != nil {
				return errors.Join(err, errors.New("error while inserting task"))
			}
			err = db.Insert(models.ConnectedTask{
				TasksId:      task.Id,
				NotionId:     notionId,
				TaskUpdate:   updated,
				ConnectionId: connectionId,
				NotionUpdate: nil,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (_ GoogleTaskService) Insert(connectionId string, details *models.TaskDetails) (string, *time.Time, error) {
	done := "needsAction"
	if details.Done {
		done = "completed"
	}
	task, err := auth.TasksService.Tasks.Insert(connectionId, &tasks.Task{
		Title:     details.Title,
		Due:       details.DueDate.Format(time.RFC3339),
		Notes:     details.Notes,
		Completed: &done,
	}).Do()
	if err != nil {
		return "", nil, errors.Join(err, errors.New("error while inserting task"))
	}
	updated, err := time.Parse(time.RFC3339, task.Updated)
	if err != nil {
		return "", nil, errors.Join(err, errors.New("error while parsing updated date"))
	}
	return task.Id, &updated, nil
}

func (_ GoogleTaskService) Update(connectionId string, id string, details *models.TaskDetails) (*time.Time, error) {
	done := "needsAction"
	if details.Done {
		done = "completed"
	}
	task, err := auth.TasksService.Tasks.Update(connectionId, id, &tasks.Task{
		Title:     details.Title,
		Due:       details.DueDate.Format(time.RFC3339),
		Notes:     details.Notes,
		Completed: &done,
	}).Do()
	if err != nil {
		return nil, errors.Join(err, errors.New("error while updating task"))
	}
	updated, err := time.Parse(time.RFC3339, task.Updated)
	if err != nil {
		return nil, errors.Join(err, errors.New("error while parsing updated date"))
	}
	return &updated, nil
}

func (_ GoogleTaskService) GetTaskDetails(connectionId string, id string) (*models.TaskDetails, error) {
	task, err := auth.TasksService.Tasks.Get(connectionId, id).Do()
	if err != nil {
		return nil, errors.Join(err, errors.New("error while getting task"))
	}
	notesArr := strings.Split(task.Notes, keys.BREAK_LINE)
	notes := notesArr[0]
	if len(notesArr) != 2 {
		notes = ""
	}
	dueDate, err := time.Parse(time.RFC3339, task.Due)
	if err != nil {
		return nil, errors.Join(err, errors.New("error while parsing due date"))
	}
	taskDetails := &models.TaskDetails{
		Title:   task.Title,
		Done:    task.Status == "completed",
		Notes:   notes,
		DueDate: &dueDate,
	}
	return taskDetails, nil
}
