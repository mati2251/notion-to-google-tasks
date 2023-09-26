package utils

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jomei/notionapi"
	"github.com/spf13/viper"
	"google.golang.org/api/tasks/v1"
)

const LAST_TIME_SYNC = "last_time_sync"

type connection struct {
	notionDatabase notionapi.DatabaseID
	tasksList      *tasks.TaskList
}

func Init() {
	var err error
	notionClient, err = GetNotionToken()
	if err != nil {
		log.Fatalf("Error getting notion client: %v", err)
	}
	tasksService, err = GetTasksService()
	if err != nil {
		log.Fatalf("Error getting google client: %v", err)
	}
}

func Sync() {
	// connections := getConnections()
	// for _, connection := range connections {
	// }
}

func ForceSync() {
	connections := getConnections()
	for _, connection := range connections {
		items, _ := notionClient.Database.Query(context.Background(), connection.notionDatabase, nil)
		for index, item := range items.Results {
			fmt.Printf("Item: %v\n", item)
			if index == 0 {
				createDbPropTasksIdIfNotExists(connection.notionDatabase)
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

func getConnections() []connection {
	var connections []connection
	for notionDatabaseId, tasksListId := range viper.GetStringMapString("connections") {
		connections = append(connections, getConnectionFromIds(notionDatabaseId, tasksListId))
	}
	return connections
}

func getConnectionFromIds(notionDatabaseId string, tasksListId string) connection {
	notionDatabaseIdObj := notionapi.DatabaseID(notionDatabaseId)
	tasksList, err := tasksService.Tasklists.Get(tasksListId).Do()
	if err != nil {
		log.Fatalf("Error getting google tasklist: %v", err)
	}
	return connection{
		notionDatabase: notionDatabaseIdObj,
		tasksList:      tasksList,
	}
}

func SetLastTimeSync() {
	viper.Set(LAST_TIME_SYNC, time.Now())
	viper.WriteConfig()
}

func GetLastTimeSync() time.Time {
	return viper.GetTime(LAST_TIME_SYNC)
}
