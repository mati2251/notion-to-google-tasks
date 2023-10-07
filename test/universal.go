package test

import (
	"context"
	"os"

	"github.com/jomei/notionapi"
	"github.com/mati2251/notion-to-google-tasks/config/auth"
	"github.com/mati2251/notion-to-google-tasks/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func InitViper() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")
	viper.SetConfigName(".notion-to-google-tasks")
	viper.AutomaticEnv()
	viper.ReadInConfig()
}

func CreateMockConnection() models.Connection {
	InitViper()
	err := auth.InitConnections()
	if err != nil {
		panic(err)
	}
	taskListId := viper.GetString("google.test_list")
	if taskListId == "" {
		panic("google.test_list is not set")
	}
	taskList, err := auth.TasksService.Tasklists.Get(taskListId).Do()
	if err != nil {
		panic(err)
	}
	notionDbId := notionapi.DatabaseID(viper.GetString("notion.test_db"))
	if notionDbId == "" {
		panic("notion.test_db is not set")
	}
	newDb, err := auth.NotionClient.Database.Get(context.Background(), notionDbId)
	if err != nil {
		panic(err)
	}
	return models.Connection{
		TasksList:      taskList,
		NotionDatabase: newDb,
	}
}

func CleanUpMock(conn models.Connection) {
	err := auth.TasksService.Tasklists.Delete(conn.TasksList.Id).Do()
	if err != nil {
		panic(err)
	}
}
