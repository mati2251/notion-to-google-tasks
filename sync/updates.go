package sync

import (
	"github.com/mati2251/notion-to-google-tasks/db"
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
	if err != nil {
		return nil, err
	}
	return ids, nil
}
