package sync

import (
	"time"

	"github.com/jomei/notionapi"
	"google.golang.org/api/tasks/v1"
)

func updateTask(notionPage notionapi.Page, googleTask *tasks.Task, force bool) error {
	googleTime, err := time.Parse(time.RFC3339, googleTask.Updated)
	if err != nil {
		return err
	}
	if notionPage.LastEditedTime.After(googleTime) || force {

	}
	return nil
}

func updateGoogle(notionPage notionapi.Page, googleTask *tasks.Task) {

}
