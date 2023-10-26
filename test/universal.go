package test

import (
	"context"
	"os"
	"time"

	"github.com/jomei/notionapi"
	"github.com/mati2251/notion-to-google-tasks/config/auth"
	"github.com/mati2251/notion-to-google-tasks/keys"
	"github.com/mati2251/notion-to-google-tasks/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func InitViper() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	configPath := home + keys.FILES_PATH
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		os.MkdirAll(configPath, 0755)
	}
	viper.AddConfigPath(configPath)
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
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
		TasksListId:     taskList.Id,
		NotionDatabasId: newDb.ID.String(),
	}
}

func CreateDetails() models.TaskDetails {
	time := time.Now().Round(time.Minute)
	return models.TaskDetails{
		Title:   "Test task",
		Notes:   "Test notes",
		Done:    false,
		DueDate: &time,
	}
}
