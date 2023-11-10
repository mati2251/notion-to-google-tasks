package test

import (
	"context"
	"path"
	"runtime"
	"time"

	"github.com/jomei/notionapi"
	"github.com/mati2251/notion-to-google-tasks/config/auth"
	"github.com/mati2251/notion-to-google-tasks/keys"
	"github.com/mati2251/notion-to-google-tasks/models"
	"github.com/spf13/viper"
	"google.golang.org/api/tasks/v1"
)

func InitViper() {
	_, b, _, _ := runtime.Caller(0)
	path := path.Join(path.Dir(path.Dir(b)))
	viper.AddConfigPath(path)
	viper.SetConfigName("test")
	viper.SetConfigType("yaml")
	error := viper.ReadInConfig()
	keys.DB_FILE = "/test.db"
	if error != nil {
		panic(error)
	}
}

func GetTestConnection() models.Connection {
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
		TasksListId:      taskList.Id,
		NotionDatabaseId: newDb.ID.String(),
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

func CreateTestTasks(conn models.Connection) (string, string) {
	taskDetails := CreateDetails()
	taskDetails.DueDate = nil
	notion, err := auth.NotionClient.Page.Create(context.Background(), &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			DatabaseID: notionapi.DatabaseID(conn.NotionDatabaseId),
			Type:       "database_id",
		},
		Properties: notionapi.Properties{
			viper.GetString(keys.NOTION_NAME_KEY): notionapi.TitleProperty{
				Title: []notionapi.RichText{{
					Type: "text",
					Text: &notionapi.Text{
						Content: taskDetails.Title,
					},
				}},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	taskDetails.Title = "Test task 2"
	taskListId := viper.GetString("google.test_list")
	task, err := auth.TasksService.Tasks.Insert(taskListId, &tasks.Task{
		Title: "Test task 2",
	}).Do()
	if err != nil {
		panic(err)
	}
	return notion.ID.String(), task.Id
}
