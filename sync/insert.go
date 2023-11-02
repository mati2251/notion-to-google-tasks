package sync

import (
	"context"
	"errors"

	"github.com/jomei/notionapi"
	"github.com/mati2251/notion-to-google-tasks/config/auth"
	"github.com/mati2251/notion-to-google-tasks/db"
	"github.com/mati2251/notion-to-google-tasks/google"
	"github.com/mati2251/notion-to-google-tasks/keys"
	"github.com/mati2251/notion-to-google-tasks/models"
	"github.com/mati2251/notion-to-google-tasks/notion"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
)

func inserts(ids []string, connectionId string) error {
	err := googleInserts(ids, connectionId)
	if err != nil {
		return err
	}
	return notionInserts(ids, connectionId)
}

func notionInserts(insertedIds []string, connectionId string) error {
	list, err := auth.TasksService.Tasks.List(connectionId).Do()
	if err != nil {
		return errors.Join(err, errors.New("error while getting tasklist"))
	}
	for _, task := range list.Items {
		if !slices.Contains(insertedIds, task.Id) {
			taskDetails, taskUpdated, err := google.Service.GetTaskDetails(connectionId, task.Id)
			if err != nil {
				return err
			}
			notionId, notionUpdated, err := notion.Service.Insert(connectionId, taskDetails)
			if err != nil {
				return errors.Join(err, errors.New("error while inserting task"))
			}
			err = db.Insert(models.ConnectedTask{
				TasksId:      task.Id,
				NotionId:     notionId,
				TaskUpdate:   taskUpdated,
				ConnectionId: connectionId,
				NotionUpdate: notionUpdated,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func googleInserts(insertedIds []string, connectionId string) error {
	notionId := notionapi.DatabaseID(viper.GetString(keys.CONNECTIONS))
	items, err := auth.NotionClient.Database.Query(context.Background(), notionId, &notionapi.DatabaseQueryRequest{})
	if err != nil {
		return errors.Join(err, errors.New("error while getting database"))
	}
	for _, page := range items.Results {
		if !slices.Contains(insertedIds, page.ID.String()) {
			details, notionUpdated, err := notion.Service.GetTaskDetails(connectionId, page.ID.String())
			if err != nil {
				return err
			}
			taskId, taskUpdated, err := google.Service.Insert(connectionId, details)
			if err != nil {
				return errors.Join(err, errors.New("error while inserting task"))
			}
			err = db.Insert(models.ConnectedTask{
				TasksId:      taskId,
				NotionId:     page.ID.String(),
				TaskUpdate:   taskUpdated,
				ConnectionId: connectionId,
				NotionUpdate: notionUpdated,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
