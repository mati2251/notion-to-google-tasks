package utils

import (
	"context"
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/jomei/notionapi"
	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
	"google.golang.org/api/tasks/v1"
)

var notionClient *notionapi.Client = nil
var tasksService *tasks.Service = nil
var bold = color.New(color.Bold)

func ConfigConnections(tasks *tasks.Service, notion *notionapi.Client) {
	notionClient = notion
	tasksService = tasks
	bold.Println("Share notion pages which you want synchronize and type ENTER")
	fmt.Scanln()
	newConnections()
}

func newConnections() {
	for {
		database, err := getNotionDb()
		if err != nil {
			log.Fatalf("Prompt failed %v\n", err)
			return
		}
		if database == nil {
			break
		}
		bold.Printf("Select Google Tasks List to synchronize with %s database\n", database.Title[0].PlainText)
		list, err := getTasksList(database.Title[0].PlainText)
		if err != nil {
			log.Fatalf("Prompt failed %v\n", err)
			return
		}
		addNewConnectionToConfig(string(database.ID), list.Id)
	}
}

func getNotionDb() (*notionapi.Database, error) {
	prompt := promptui.Prompt{
		Label: "Type notion database id (type ENTER to exit)",
	}
	databaseIdString, err := prompt.Run()
	if err != nil {
		return nil, err
	}
	if databaseIdString == "" {
		return nil, nil
	}

	databaseId := notionapi.DatabaseID(databaseIdString)
	database, err := notionClient.Database.Get(context.Background(), databaseId)
	if err != nil {
		return nil, err
	}
	return database, err
}

func getTasksList(eventuallyNewName string) (*tasks.TaskList, error) {
	lists, _ := tasksService.Tasklists.List().Do()
	var taskListTitles []string
	for _, list := range lists.Items {
		taskListTitles = append(taskListTitles, list.Title)
	}
	newListKey := "New list"
	taskListTitles = append(taskListTitles, newListKey)
	prompt := promptui.Select{
		Label: "Select task list",
		Items: taskListTitles,
	}
	index, result, err := prompt.Run()
	if err != nil {
		return nil, err
	}
	if result == newListKey {
		return createNewList(eventuallyNewName)
	}
	return lists.Items[index], nil
}

func addNewConnectionToConfig(databaseId string, listId string) {
	viper.Set(fmt.Sprintf("connections.%s", databaseId), listId)
}

func createNewList(name string) (*tasks.TaskList, error) {
	newList := &tasks.TaskList{
		Title: name,
	}
	return tasksService.Tasklists.Insert(newList).Do()
}