package sync

import (
	"errors"
	"log"

	"github.com/mati2251/notion-to-google-tasks/db"
	"github.com/mati2251/notion-to-google-tasks/google"
	"github.com/mati2251/notion-to-google-tasks/models"
	"github.com/mati2251/notion-to-google-tasks/notion"
)

func updates(connectionId string) ([]string, error) {
	connectedTasks, err := db.GetConnectedTasks(connectionId)
	ids := []string{}
	if err != nil {
		return nil, err
	}
	for _, connectedTask := range connectedTasks {
		ids = append(ids, connectedTask.TasksId)
		ids = append(ids, connectedTask.NotionId)
		err := update(connectedTask)
		if err != nil {
			log.Println(err)
		}
	}
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func update(task models.ConnectedTask) error {
	notionTask, notionUpdated, err := notion.Service.GetTaskDetails(task.ConnectionId, task.NotionId)
	if err != nil {
		return errors.Join(errors.New("error during updated task"), err)
	}
	if notionUpdated.After(*task.NotionUpdate) {
		googleUpdated, err := google.Service.Update(task.ConnectionId, task.TasksId, notionTask)
		if err != nil {
			return errors.Join(errors.New("error during updated google task"), err)
		}
		err = db.UpdateNotionTime(task.NotionId, notionUpdated)
		if err != nil {
			return errors.Join(errors.New("error during updated time task in db"), err)
		}
		err = db.UpdateGoogleTime(task.TasksId, googleUpdated)
		if err != nil {
			return errors.Join(errors.New("error during updated time task in db"), err)
		}
	} else if googleTask, googleUpdated, err := google.Service.GetTaskDetails(task.ConnectionId, task.TasksId); err == nil {
		if googleUpdated.After(*task.TaskUpdate) {
			notionUpdated, err := notion.Service.Update(task.ConnectionId, task.NotionId, googleTask)
			if err != nil {
				return errors.Join(errors.New("error during updated notion task"), err)
			}
			err = db.UpdateNotionTime(task.NotionId, notionUpdated)
			if err != nil {
				return errors.Join(errors.New("error during updated time task in db"), err)
			}
			err = db.UpdateGoogleTime(task.TasksId, googleUpdated)
			if err != nil {
				return errors.Join(errors.New("error during updated time task in db"), err)
			}
		}
	}
	return nil
}
