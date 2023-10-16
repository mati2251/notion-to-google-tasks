package sync

import (
	"github.com/mati2251/notion-to-google-tasks/db"
	"github.com/mati2251/notion-to-google-tasks/models"
)

func Sync(connections []models.Connection) error {
	for _, connection := range connections {
		err := db.OpenFile()
		if err != nil {
			return err
		}
		ids, err := updates()
		err = inserts(ids, connection.TasksListId)
		if err != nil {
			return err
		}
	}
	return nil
}
