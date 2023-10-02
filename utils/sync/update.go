package sync

import (
	"github.com/jomei/notionapi"
	"google.golang.org/api/tasks/v1"
)

func updateConnections(notionPage notionapi.Page, googleTask *tasks.Task) {
	googleTime := googleTask.Updated
	println(googleTime)
	// if notionPage.LastEditedTime.After() {
	// 	// todo update google task
	// } else {
	// 	// todo update notion page
	// }
}
