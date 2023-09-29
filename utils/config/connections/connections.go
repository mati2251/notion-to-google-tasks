package connections

import (
	"log"

	"github.com/jomei/notionapi"
	"github.com/mati2251/notion-to-google-tasks/utils/config/auth"
	"github.com/spf13/viper"
	"google.golang.org/api/tasks/v1"
)

type Connection struct {
	NotionDatabase notionapi.DatabaseID
	TasksList      *tasks.TaskList
}

func GetConnections() []Connection {
	var connections []Connection
	for notionDatabaseId, tasksListId := range viper.GetStringMapString("connections") {
		connections = append(connections, getConnectionFromIds(notionDatabaseId, tasksListId))
	}
	return connections
}

func getConnectionFromIds(notionDatabaseId string, tasksListId string) Connection {
	notionDatabaseIdObj := notionapi.DatabaseID(notionDatabaseId)
	tasksList, err := auth.TasksService.Tasklists.Get(tasksListId).Do()
	if err != nil {
		log.Fatalf("Error getting google tasklist: %v", err)
	}
	return Connection{
		NotionDatabase: notionDatabaseIdObj,
		TasksList:      tasksList,
	}
}
