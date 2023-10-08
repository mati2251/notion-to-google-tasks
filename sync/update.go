package sync

import (
	"errors"
	"time"

	"github.com/mati2251/notion-to-google-tasks/models"
	"github.com/mati2251/notion-to-google-tasks/sync/google"
	"github.com/mati2251/notion-to-google-tasks/sync/notion"
)

func update(connectedTask models.ConnectedTask, force bool) (models.ConnectedTask, error) {
	if connectedTask.Task == nil {
		page, err := notion.SetDone(connectedTask.Notion)
		if err != nil {
			return connectedTask, errors.Join(err, errors.New("error setting done"))
		}
		connectedTask.Notion = page
		return connectedTask, nil
	}
	googleTime, err := time.Parse(time.RFC3339, connectedTask.Task.Updated)
	googleTime = googleTime.Add(-time.Duration(googleTime.Second()) * time.Second)
	if err != nil {
		return connectedTask, err
	}
	if connectedTask.Notion == nil {
		task, err := google.SetDone(*connectedTask.Connection.TasksList, connectedTask.Task)
		if err != nil {
			return connectedTask, errors.Join(err, errors.New("error setting done"))
		}
		connectedTask.Task = task
		return connectedTask, nil
	}
	if connectedTask.Notion.LastEditedTime.Before(googleTime) {
		return notion.Update(connectedTask)
	} else if connectedTask.Notion.LastEditedTime.After(googleTime) || force {
		return google.Update(connectedTask)
	}
	return connectedTask, nil
}
