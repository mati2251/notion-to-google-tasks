package sync

import (
	"github.com/mati2251/notion-to-google-tasks/db"
	"github.com/mati2251/notion-to-google-tasks/models"
)

func Sync(connections []models.Connection) error {
	err := db.OpenFile()
	if err != nil {
		return err
	}
	for _, connection := range connections {
		ids, err := updates(connection.NotionDatabaseId)
		if err != nil {
			return err
		}
		err = inserts(ids, connection.NotionDatabaseId)
		if err != nil {
			return err
		}
	}
	return nil
}
