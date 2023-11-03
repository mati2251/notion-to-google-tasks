package sync

import (
	"github.com/mati2251/notion-to-google-tasks/db"
	"github.com/mati2251/notion-to-google-tasks/models"
)

func updates(connectionId string) ([]string, error) {
	connectedTasks, err := db.GetConnectedTasks(connectionId)
	ids := []string{}
	if err != nil {
		return nil, err
	}
	for _, connectedTask := range connectedTasks {
		ids = append(ids, connectedTask.TasksId)
		ids = append(ids, connectedTask.NotionId)
	}
	err = inserts(ids, connectionId)
	if err != nil {
		return nil, err
	}
	return []string{}, nil
}

func update(connectedTask models.ConnectedTask) {
	
}