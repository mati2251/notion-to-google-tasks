package utils

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jomei/notionapi"
	"github.com/spf13/viper"
	"google.golang.org/api/tasks/v1"
)

const LAST_TIME_SYNC = "last_time_sync"

type connection struct {
	notionDatabase notionapi.DatabaseID
	tasksList      *tasks.TaskList
}

func Init() {
	var err error
	notionClient, err = GetNotionToken()
	if err != nil {
		log.Fatalf("Error getting notion client: %v", err)
	}
	tasksService, err = GetTasksService()
	if err != nil {
		log.Fatalf("Error getting google client: %v", err)
	}
}

func Sync() {
	// connections := getConnections()
	// for _, connection := range connections {
	// }
}

func ForceSync() {
	connections := getConnections()
	for _, connection := range connections {
		items, _ := notionClient.Database.Query(context.Background(), connection.notionDatabase, nil)
		for index, item := range items.Results {
			if index == 0 {
				createDbPropTasksIdIfNotExists(connection.notionDatabase)
			}
			tasksId := getStringValueFromProperty(item.Properties["Tasks ID"])
			if tasksId == "" {
				nameKey := viper.GetString(NOTION_NAME_KEY)
				if item.Properties[nameKey] == nil {
					log.Fatalf("Invalid notion name key: %v", nameKey)
				}
				name := getStringValueFromProperty(item.Properties[nameKey])
				taskId := createNewTaskAtGoogle(name, connection)
			} else {
				// todo update tasks
			}
		}
	}
}

func createDbPropTasksIdIfNotExists(databaseId notionapi.DatabaseID) {
	result, _ := notionClient.Database.Query(context.Background(), databaseId, nil)
	if result.Results[0].Properties["Tasks ID"] == nil {
		createDbPropTasksId(databaseId)
	}
}

func createNewTaskAtGoogle(name string, connection connection) string {
	newTask := &tasks.Task{
		Title: name,
	}
	task, err := tasksService.Tasks.Insert(connection.tasksList.Id, newTask).Do()
	if err != nil {
		log.Fatalf("Error creating task: %v", err)
	}
	return task.Id
}

func getStringValueFromProperty(property notionapi.Property) string {
	switch property.GetType() {
	case notionapi.PropertyTypeRichText:
		richText := property.(*notionapi.RichTextProperty).RichText
		var value string
		for index, richTextItem := range richText {
			value += richTextItem.PlainText
			if index != len(richText)-1 {
				value += " "
			}
		}
		return value
	case notionapi.PropertyTypeText:
		textProperty := property.(*notionapi.TextProperty).Text
		var value string
		for index, item := range textProperty {
			value += item.PlainText
			if index != len(textProperty)-1 {
				value += " "
			}
		}
		return value
	case notionapi.PropertyTypeTitle:
		textProperty := property.(*notionapi.TitleProperty).Title
		var value string
		for index, item := range textProperty {
			value += item.PlainText
			if index != len(textProperty)-1 {
				value += " "
			}
		}
		return value
	case notionapi.PropertyTypeNumber:
		return fmt.Sprintf("%v", property.(*notionapi.NumberProperty).Number)
	case notionapi.PropertyTypeSelect:
		return property.(*notionapi.SelectProperty).Select.Name
	case notionapi.PropertyTypeMultiSelect:
		multiSelect := property.(*notionapi.MultiSelectProperty).MultiSelect
		var value string
		for index, item := range multiSelect {
			value += item.Name
			if index != len(multiSelect)-1 {
				value += " "
			}
		}
		return value
	case notionapi.PropertyTypeDate:
		return property.(*notionapi.DateProperty).Date.Start.String()
	case notionapi.PropertyTypeFormula:
		return property.(*notionapi.FormulaProperty).Formula.String
	case notionapi.PropertyTypeRelation:
		relation := property.(*notionapi.RelationProperty).Relation
		var value string
		for index, item := range relation {
			value += item.ID.String() + " "
			if index != len(relation)-1 {
				value += " "
			}
		}
		return value
	case notionapi.PropertyTypeRollup:
		rollup := property.(*notionapi.RollupProperty).Rollup
		return rollup.Date.End.String() + " " + rollup.Date.Start.String()
	case notionapi.PropertyTypePeople:
		people := property.(*notionapi.PeopleProperty).People
		var value string
		for index, item := range people {
			value += item.Name + " "
			if index != len(people)-1 {
				value += " "
			}
		}
		return value
	case notionapi.PropertyTypeFiles:
		files := property.(*notionapi.FilesProperty).Files
		var value string
		for index, item := range files {
			value += item.Name + " "
			if index != len(files)-1 {
				value += " "
			}
		}
		return value
	case notionapi.PropertyTypeCheckbox:
		return fmt.Sprintf("%v", property.(*notionapi.CheckboxProperty).Checkbox)
	case notionapi.PropertyTypeURL:
		return property.(*notionapi.URLProperty).URL
	case notionapi.PropertyTypeEmail:
		return property.(*notionapi.EmailProperty).Email
	case notionapi.PropertyTypePhoneNumber:
		return property.(*notionapi.PhoneNumberProperty).PhoneNumber
	case notionapi.PropertyTypeCreatedTime:
		return property.(*notionapi.CreatedTimeProperty).CreatedTime.String()
	case notionapi.PropertyTypeCreatedBy:
		return property.(*notionapi.CreatedByProperty).CreatedBy.Name
	case notionapi.PropertyTypeLastEditedTime:
		return property.(*notionapi.LastEditedTimeProperty).LastEditedTime.String()
	case notionapi.PropertyTypeLastEditedBy:
		return property.(*notionapi.LastEditedByProperty).LastEditedBy.Name
	case notionapi.PropertyTypeStatus:
		return property.(*notionapi.StatusProperty).Status.Name
	case notionapi.PropertyTypeUniqueID:
		return property.(*notionapi.UniqueIDProperty).UniqueID.String()
	case notionapi.PropertyTypeVerification:
		return property.(*notionapi.VerificationProperty).Verification.VerifiedBy.Name
	default:
		return ""
	}
}

func createDbPropTasksId(databaseId notionapi.DatabaseID) {
	defaultValue := &notionapi.RichTextPropertyConfig{
		Type: "rich_text",
	}
	properties := notionapi.PropertyConfigs(map[string]notionapi.PropertyConfig{"Tasks ID": defaultValue})
	_, err := notionClient.Database.Update(context.Background(), databaseId, &notionapi.DatabaseUpdateRequest{
		Properties: properties,
	})
	if err != nil {
		log.Fatalf("Error updating database: %v", err)
	}
}

func getConnections() []connection {
	var connections []connection
	for notionDatabaseId, tasksListId := range viper.GetStringMapString("connections") {
		connections = append(connections, getConnectionFromIds(notionDatabaseId, tasksListId))
	}
	return connections
}

func getConnectionFromIds(notionDatabaseId string, tasksListId string) connection {
	notionDatabaseIdObj := notionapi.DatabaseID(notionDatabaseId)
	tasksList, err := tasksService.Tasklists.Get(tasksListId).Do()
	if err != nil {
		log.Fatalf("Error getting google tasklist: %v", err)
	}
	return connection{
		notionDatabase: notionDatabaseIdObj,
		tasksList:      tasksList,
	}
}

func SetLastTimeSync() {
	viper.Set(LAST_TIME_SYNC, time.Now())
	viper.WriteConfig()
}

func GetLastTimeSync() time.Time {
	return viper.GetTime(LAST_TIME_SYNC)
}
