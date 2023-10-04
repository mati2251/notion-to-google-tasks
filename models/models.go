package models

import (
	"context"
	"log"

	"github.com/jomei/notionapi"
	"github.com/mati2251/notion-to-google-tasks/config/auth"
	"google.golang.org/api/tasks/v1"
)

type ConnectedTask struct {
	Notion     *notionapi.Page
	Task       *tasks.Task
	Connection *Connection
}

type Connection struct {
	NotionDatabase *notionapi.Database
	TasksList      *tasks.TaskList
}

func GetConnectionFromIds(notionDatabaseId string, tasksListId string) Connection {
	notionDatabaseIdObject := notionapi.DatabaseID(notionDatabaseId)
	notionDatabase, err := auth.NotionClient.Database.Get(context.Background(), notionDatabaseIdObject)
	tasksList, err := auth.TasksService.Tasklists.Get(tasksListId).Do()
	if err != nil {
		log.Fatalf("Error getting google tasklist: %v", err)
	}
	if err != nil {
		log.Fatalf("Error getting notion database: %v", err)
	}
	return Connection{
		NotionDatabase: notionDatabase,
		TasksList:      tasksList,
	}
}
