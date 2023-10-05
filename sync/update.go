package sync

import (
	"time"

	"github.com/mati2251/notion-to-google-tasks/models"
	"github.com/mati2251/notion-to-google-tasks/sync/google"
	"github.com/mati2251/notion-to-google-tasks/sync/notion"
)

func updateTask(connectedTask models.ConnectedTask, force bool) error {
	googleTime, err := time.Parse(time.RFC3339, connectedTask.Task.Updated)
	if err != nil {
		return err
	}
	if connectedTask.Notion.LastEditedTime.Before(googleTime) {
		notion.Update(connectedTask)
	} else if connectedTask.Notion.LastEditedTime.After(googleTime) || force {
		google.Update(connectedTask)
	}
	return nil
}
