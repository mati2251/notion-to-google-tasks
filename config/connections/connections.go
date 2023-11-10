package connections

import (
	"github.com/mati2251/notion-to-google-tasks/keys"
	"github.com/mati2251/notion-to-google-tasks/models"
	"github.com/spf13/viper"
)

func GetConnections() []models.Connection {
	var connections []models.Connection
	for notionDatabaseId, tasksListId := range viper.GetStringMapString(keys.CONNECTIONS) {
		connections = append(connections, models.Connection{NotionDatabaseId: notionDatabaseId, TasksListId: tasksListId})
	}
	return connections
}
