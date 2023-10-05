package test

import (
	"context"
	"os"

	"github.com/jomei/notionapi"
	"github.com/mati2251/notion-to-google-tasks/config/auth"
	"github.com/mati2251/notion-to-google-tasks/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/api/tasks/v1"
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
	new := &tasks.TaskList{Title: "test"}
	newTaskList, err := auth.TasksService.Tasklists.Insert(new).Do()
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
		TasksList:      newTaskList,
		NotionDatabase: newDb,
	}
}

func CreateMockPage(conn models.Connection) {
	// notionDatabaseId := notionapi.DatabaseID(conn.NotionDatabase.ID)
	// page, err := auth.NotionClient.Page.Create(context.Background(), &notionapi.PageCreateRequest{
	// 	Parent: notionapi.Parent{
	// 		Type:       "database_id",
	// 		DatabaseID: notionDatabaseId,
	// 		// Properties: notionapi.Properties
	// 	}})
}

func CreateMockTask() {

}
