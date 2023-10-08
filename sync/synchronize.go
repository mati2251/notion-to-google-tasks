package sync

import (
	"context"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/jomei/notionapi"
	"github.com/mati2251/notion-to-google-tasks/config/auth"
	"github.com/mati2251/notion-to-google-tasks/config/connections"
	"github.com/mati2251/notion-to-google-tasks/keys"
	"github.com/mati2251/notion-to-google-tasks/models"
	"github.com/mati2251/notion-to-google-tasks/sync/google"
	"github.com/mati2251/notion-to-google-tasks/sync/notion"
	"github.com/spf13/viper"
)

func Sync(force bool) {
	connections := connections.GetConnections()
	for _, connection := range connections {
		check, err := checkIsSyncRequired(connection, force)
		if err != nil {
			log.Fatalf("Error checking connection: %v", err)
		}
		if check {
			itemsId := syncFromNotion(connection, force)
			syncFromGoogle(connection, force, itemsId)
			SetLastTimeSync()
		}
	}
}

func checkIsSyncRequired(connection models.Connection, force bool) (bool, error) {
	if force {
		return true, nil
	}
	googleTime, err := time.Parse(time.RFC3339, connection.TasksList.Updated)
	if err != nil {
		return false, err
	}
	googleTime = googleTime.Add(-time.Duration(googleTime.Second()) * time.Second)
	notionTime := connection.NotionDatabase.LastEditedTime
	lastTimeSync := GetLastTimeSync()
	return lastTimeSync.Before(googleTime) || lastTimeSync.Before(notionTime), nil
}

func syncFromNotion(connection models.Connection, force bool) []string {
	itemsId := make([]string, 0)
	databaseId := notionapi.DatabaseID(connection.NotionDatabase.ID)
	items, _ := auth.NotionClient.Database.Query(context.Background(), databaseId, nil)
	notion.CreateProp(*connection.NotionDatabase, keys.TASK_ID_KEY, "rich_text")
	for _, item := range items.Results {
		itemsId = append(itemsId, item.ID.String())
		tasksId := notion.GetStringValueFromProperty(item.Properties[keys.TASK_ID_KEY])
		if tasksId == "" {
			google.New(models.ConnectedTask{
				Notion:     &item,
				Task:       nil,
				Connection: &connection,
			})
		} else {
			task, err := auth.TasksService.Tasks.Get(connection.TasksList.Id, tasksId).Do()
			connectedTask := models.ConnectedTask{
				Notion:     &item,
				Task:       task,
				Connection: &connection,
			}
			if err != nil {
				notion.UpdateValueFromProp(&item, viper.GetString(keys.NOTION_STATUS_KEY), viper.GetString(keys.NOTION_DONE_STATUS_VALUE))
			} else {
				update(connectedTask, force)
			}
		}
	}
	return itemsId
}

func syncFromGoogle(connection models.Connection, force bool, itemsId []string) {
	tasks, _ := auth.TasksService.Tasks.List(connection.TasksList.Id).Do()
	for _, task := range tasks.Items {
		if !slices.Contains(itemsId, task.Id) {
			if strings.Contains(task.Notes, "https://www.notion.so/") {
				_, err := google.SetDone(*connection.TasksList, task)
				if err != nil {
					log.Fatalf("Error setting task done: %v", err)
				}
			} else {
				notion.New(models.ConnectedTask{
					Notion:     nil,
					Task:       task,
					Connection: &connection,
				})
			}
		}
	}
}
