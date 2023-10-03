package models

import (
	"log"

	"github.com/jomei/notionapi"
	"github.com/mati2251/notion-to-google-tasks/utils/config/auth"
	"google.golang.org/api/tasks/v1"
)

type ConnectedTask struct {
	Notion     *notionapi.Page
	Task       *tasks.Task
	Connection *Connection
}

type Connection struct {
	NotionDatabase notionapi.DatabaseID
	TasksList      *tasks.TaskList
}

func GetConnectionFromIds(notionDatabaseId string, tasksListId string) Connection {
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
