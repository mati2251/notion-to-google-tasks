package auth

import (
	"github.com/jomei/notionapi"
	"google.golang.org/api/tasks/v1"
)

var NotionClient *notionapi.Client = nil
var TasksService *tasks.Service = nil

func InitConnections() error {
	var notionErr, tasksErr error
	NotionClient, notionErr = GetNotionToken()
	TasksService, tasksErr = GetTasksService()
	if notionErr != nil {
		return notionErr
	}
	if tasksErr != nil {
		return tasksErr
	}
	return nil
}
