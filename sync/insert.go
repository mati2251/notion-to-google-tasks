package sync

import (
	"github.com/mati2251/notion-to-google-tasks/google"
	"github.com/mati2251/notion-to-google-tasks/notion"
)

func inserts(ids []string, connectionId string) error {
	err := notion.Service.Inserts(ids, connectionId)
	if err != nil {
		return err
	}
	return google.Service.Inserts(ids, connectionId)
}
