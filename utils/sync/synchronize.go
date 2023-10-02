package sync

import (
	"context"

	"github.com/mati2251/notion-to-google-tasks/utils/config/auth"
	"github.com/mati2251/notion-to-google-tasks/utils/config/connections"
	"github.com/mati2251/notion-to-google-tasks/utils/keys"
)

func Sync() {
	// connections := getConnections()
	// for _, connection := range connections {
	// }
}

func ForceSync() {
	connections := connections.GetConnections()
	for _, connection := range connections {
		items, _ := auth.NotionClient.Database.Query(context.Background(), connection.NotionDatabase, nil)
		createDbPropTasksIdIfNotExists(connection.NotionDatabase)
		for _, item := range items.Results {
			tasksId := getStringValueFromProperty(item.Properties[keys.TASK_ID_KEY])
			if tasksId == "" {
				createNewTaskAtGoogle(item, connection)
			} else {
				task, err := auth.TasksService.Tasks.Get(connection.TasksList.Id, tasksId).Do()
				if err != nil {
					// check done notion task
				}
				updateConnections(item, task)
			}
		}
	}
}
