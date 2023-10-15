package google

import (
	"errors"
	"slices"

	"github.com/mati2251/notion-to-google-tasks/config/auth"
)

type GoogleTaskService struct{}

var Service GoogleTaskService

func (_ GoogleTaskService) Inserts(ids []string, connectionId string) error {
	list, err := auth.TasksService.Tasks.List(connectionId).Do()
	if err != nil {
		return errors.Join(err, errors.New("error while getting tasklist"))
	}
	for _, task := range list.Items {
		if !slices.Contains(ids, task.Id) {

		}
	}
	return nil
}

func (_ GoogleTaskService) Insert(taskId string, connectionId string) error {
	return nil
}
