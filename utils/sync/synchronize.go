package sync

import (
	"context"
	"fmt"
	"log"

	"github.com/jomei/notionapi"
	"github.com/mati2251/notion-to-google-tasks/utils/config/connections"
	"github.com/mati2251/notion-to-google-tasks/utils/keys"
	"github.com/spf13/viper"
	"google.golang.org/api/tasks/v1"
)

var notionClient *notionapi.Client
var tasksService *tasks.Service

func Sync() {
	// connections := getConnections()
	// for _, connection := range connections {
	// }
}

func ForceSync() {
	connections := connections.GetConnections()
	for _, connection := range connections {
		items, _ := notionClient.Database.Query(context.Background(), connection.NotionDatabase, nil)
		createDbPropTasksIdIfNotExists(connection.NotionDatabase)
		for _, item := range items.Results {
			tasksId := getStringValueFromProperty(item.Properties["Tasks ID"])
			if tasksId == "" {
				createNewTaskAtGoogle(item, connection)
			} else {
				// todo update tasks
			}
		}
	}
}

func createDbPropTasksIdIfNotExists(databaseId notionapi.DatabaseID) {
	result, _ := notionClient.Database.Query(context.Background(), databaseId, nil)
	if result.Results[0].Properties["Tasks ID"] == nil {
		createDbPropTasksId(databaseId)
	}
}

func createNewTaskAtGoogle(tuple notionapi.Page, connection connections.Connection) string {
	nameKey := viper.GetString(keys.NOTION_NAME_KEY)
	if tuple.Properties[nameKey] == nil {
		log.Fatalf("Invalid notion name key: %v", nameKey)
	}
	// name := getStringValueFromProperty(tuple.Properties[nameKey])
	for test, prop := range tuple.Properties {
		fmt.Printf("test: %v\n", test)
		fmt.Printf("prop: %v\n", prop)
	}
	// newTask := &tasks.Task{
	// Title: name,
	// }
	// task, err := tasksService.Tasks.Insert(connection.TasksList.Id, newTask).Do()
	// if err != nil {
	// log.Fatalf("Error creating task: %v", err)
	// }
	// insertTaskIdToNotion(task.Id, tuple.ID)
	// return task.Id
	return ""
}

func insertTaskIdToNotion(taskId string, notionId notionapi.ObjectID) {
	_, err := notionClient.Page.Update(context.Background(), notionapi.PageID(notionId), &notionapi.PageUpdateRequest{
		Properties: notionapi.Properties{
			"Tasks ID": &notionapi.RichTextProperty{
				RichText: []notionapi.RichText{
					{
						Type: "text",
						Text: &notionapi.Text{
							Content: taskId,
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("Error updating page: %v", err)
	}
}

func createDbPropTasksId(databaseId notionapi.DatabaseID) {
	defaultValue := &notionapi.RichTextPropertyConfig{
		Type: "rich_text",
	}
	properties := notionapi.PropertyConfigs(map[string]notionapi.PropertyConfig{"Tasks ID": defaultValue})
	_, err := notionClient.Database.Update(context.Background(), databaseId, &notionapi.DatabaseUpdateRequest{
		Properties: properties,
	})
	if err != nil {
		log.Fatalf("Error updating database: %v", err)
	}
}
