package sync

import (
	"context"
	"fmt"
	"log"

	"github.com/jomei/notionapi"
	"github.com/mati2251/notion-to-google-tasks/utils/config/auth"
	"github.com/mati2251/notion-to-google-tasks/utils/config/connections"
	"github.com/mati2251/notion-to-google-tasks/utils/keys"
	"github.com/spf13/viper"
	"google.golang.org/api/tasks/v1"
)

func createDbPropTasksIdIfNotExists(databaseId notionapi.DatabaseID) {
	result, _ := auth.NotionClient.Database.Query(context.Background(), databaseId, nil)
	if result.Results[0].Properties[keys.TASK_ID_KEY] == nil {
		createDbPropTasksId(databaseId)
	}
}

func createNewTaskAtGoogle(tuple notionapi.Page, connection connections.Connection) string {
	newTask := &tasks.Task{
		Title: getNameForTasks(tuple),
		Notes: createNotes(tuple),
		Due:   getDeadlineForTasks(tuple),
	}
	task, err := auth.TasksService.Tasks.Insert(connection.TasksList.Id, newTask).Do()
	if err != nil {
		log.Fatalf("Error creating task: %v", err)
	}
	insertTaskIdToNotion(task.Id, tuple.ID)
	return task.Id
}

func getNameForTasks(tuple notionapi.Page) string {
	nameKey := viper.GetString(keys.NOTION_NAME_KEY)
	if tuple.Properties[nameKey] == nil {
		log.Fatalf("Invalid notion name key: %v", nameKey)
	}
	return getStringValueFromProperty(tuple.Properties[nameKey])
}

func getDeadlineForTasks(tuple notionapi.Page) string {
	deadlineKey := viper.GetString(keys.NOTION_DEADLINE_KEY)
	deadline := ""
	if tuple.Properties[deadlineKey] != nil {
		deadline = getStringValueFromProperty(tuple.Properties[deadlineKey])
	}
	return deadline
}

func createNotes(tuple notionapi.Page) string {
	notes := keys.BREAK_LINE
	notes += notionPropsToString(tuple)
	notes += fmt.Sprintf("%v: %v\n", keys.NOTION_LINK_KEY, tuple.URL)
	notes += fmt.Sprintf("%v: %v\n", keys.NOTION_ID_KEY, tuple.ID)
	return notes
}

func notionPropsToString(tuple notionapi.Page) string {
	deadlineKey := viper.GetString(keys.NOTION_DEADLINE_KEY)
	nameKey := viper.GetString(keys.NOTION_NAME_KEY)
	propsString := ""
	for key, value := range tuple.Properties {
		if key != keys.TASK_ID_KEY && key != nameKey && key != deadlineKey {
			propsString += fmt.Sprintf("%v: %v\n", key, getStringValueFromProperty(value))
		}
	}
	return propsString
}

func insertTaskIdToNotion(taskId string, notionId notionapi.ObjectID) {
	_, err := auth.NotionClient.Page.Update(context.Background(), notionapi.PageID(notionId), &notionapi.PageUpdateRequest{
		Properties: notionapi.Properties{
			keys.TASK_ID_KEY: &notionapi.RichTextProperty{
				RichText: []notionapi.RichText{
					{
						Type: "text",
						Text: &notionapi.Text{
							Content: taskId,
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("Error updating page: %v", err)
	}
}

func createDbPropTasksId(databaseId notionapi.DatabaseID) {
	defaultValue := &notionapi.RichTextPropertyConfig{
		Type: "rich_text",
	}
	properties := notionapi.PropertyConfigs(map[string]notionapi.PropertyConfig{keys.TASK_ID_KEY: defaultValue})
	_, err := auth.NotionClient.Database.Update(context.Background(), databaseId, &notionapi.DatabaseUpdateRequest{
		Properties: properties,
	})
	if err != nil {
		log.Fatalf("Error updating database: %v", err)
	}
}
