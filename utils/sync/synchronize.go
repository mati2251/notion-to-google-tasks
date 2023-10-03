package sync

import (
	"context"
	"log"

	"github.com/mati2251/notion-to-google-tasks/utils/config/auth"
	"github.com/mati2251/notion-to-google-tasks/utils/config/connections"
	"github.com/mati2251/notion-to-google-tasks/utils/keys"
	"github.com/mati2251/notion-to-google-tasks/utils/models"
	"github.com/mati2251/notion-to-google-tasks/utils/sync/google"
	"github.com/mati2251/notion-to-google-tasks/utils/sync/notion"
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
		notion.CreateDbPropTasksIdIfNotExists(connection.NotionDatabase)
		for _, item := range items.Results {
			tasksId := notion.GetStringValueFromProperty(item.Properties[keys.TASK_ID_KEY])
			if tasksId == "" {
				google.CreateNewTask(models.ConnectedTask{
					Notion:     &item,
					Task:       nil,
					Connection: &connection,
				})
			} else {
				task, err := auth.TasksService.Tasks.Get(connection.TasksList.Id, tasksId).Do()
				log.Default().Println(task)
				if err != nil {
					// check done notion task
				}
				// updateTask(item, task)
			}
		}
	}
}
