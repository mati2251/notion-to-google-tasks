package connections

import (
	"github.com/mati2251/notion-to-google-tasks/models"
	"github.com/spf13/viper"
)

func GetConnections() []models.Connection {
	var connections []models.Connection
	for notionDatabaseId, tasksListId := range viper.GetStringMapString("connections") {
		connections = append(connections, models.GetConnectionFromIds(notionDatabaseId, tasksListId))
	}
	return connections
}
