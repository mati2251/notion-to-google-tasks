package notion

import (
	"context"
	"testing"
	"time"

	"github.com/jomei/notionapi"
	"github.com/mati2251/notion-to-google-tasks/config/auth"
	"github.com/mati2251/notion-to-google-tasks/keys"
	"github.com/mati2251/notion-to-google-tasks/models"
	"github.com/mati2251/notion-to-google-tasks/test"
	"github.com/spf13/viper"
)

var connection models.Connection = test.GetTestConnection()

func TestInsert(t *testing.T) {
	details := test.CreateDetails()
	id, updated, err := Service.Insert(connection.NotionDatabaseId, &details)
	if err != nil {
		t.Error(err)
	}
	page, err := auth.NotionClient.Page.Get(context.Background(), notionapi.PageID(id))
	if err != nil {
		t.Error(err)
	}
	newTitle := GetStringValueFromProperty(page.Properties[viper.GetString(keys.NOTION_NAME_KEY)])
	if newTitle != details.Title {
		t.Errorf("Title is not correct: new value %v correct value %v", newTitle, details.Title)
	}
	newDate := GetStringValueFromProperty(page.Properties[viper.GetString(keys.NOTION_DEADLINE_KEY)])
	if newDate != details.DueDate.Format(time.RFC3339) {
		t.Errorf("Deadline is not correct: new value %v correct value %v", newDate, details.DueDate.Format(time.RFC3339))
	}
	if page.LastEditedTime != *updated {
		t.Error("Last edited time is not correct")
	}
	t.Cleanup(func() {
		_, err := auth.NotionClient.Page.Update(context.Background(), notionapi.PageID(id), &notionapi.PageUpdateRequest{
			Archived:   true,
			Properties: notionapi.Properties{},
		})
		if err != nil {
			t.Error(err)
		}
	})
}

func TestInsertWithNullDate(t *testing.T) {
	details := test.CreateDetails()
	details.DueDate = nil
	id, updated, err := Service.Insert(connection.NotionDatabaseId, &details)
	if err != nil {
		t.Error(err)
	}
	page, err := auth.NotionClient.Page.Get(context.Background(), notionapi.PageID(id))
	if err != nil {
		t.Error(err)
	}
	newTitle := GetStringValueFromProperty(page.Properties[viper.GetString(keys.NOTION_NAME_KEY)])
	if newTitle != details.Title {
		t.Errorf("Title is not correct: new value %v correct value %v", newTitle, details.Title)
	}
	newDate := GetStringValueFromProperty(page.Properties[viper.GetString(keys.NOTION_DEADLINE_KEY)])
	if newDate != "" {
		t.Errorf("Deadline is not correct: new value %v correct value %v", newDate, details.DueDate.Format(time.RFC3339))
	}
	if page.LastEditedTime != *updated {
		t.Error("Last edited time is not correct")
	}
	t.Cleanup(func() {
		_, err := auth.NotionClient.Page.Update(context.Background(), notionapi.PageID(id), &notionapi.PageUpdateRequest{
			Archived:   true,
			Properties: notionapi.Properties{},
		})
		if err != nil {
			t.Error(err)
		}
	})
}

func TestGetTaskDetails(t *testing.T) {
	details := test.CreateDetails()
	pageId, _, err := Service.Insert(connection.NotionDatabaseId, &details)
	if err != nil {
		t.Error(err)
	}
	taskDetails, _, err := Service.GetTaskDetails(connection.NotionDatabaseId, pageId)
	if err != nil {
		t.Error(err)
	}
	if taskDetails.Title != details.Title {
		t.Errorf("Title is not correct: new value %v correct value %v", taskDetails.Title, details.Title)
	}
	if taskDetails.DueDate.Format(time.RFC3339) != details.DueDate.Format(time.RFC3339) {
		t.Errorf("Deadline is not correct: new value %v correct value %v", taskDetails.DueDate.Format(time.RFC3339), details.DueDate.Format(time.RFC3339))
	}
	if taskDetails.Done != details.Done {
		t.Errorf("Done is not correct: new value %v correct value %v", taskDetails.Done, details.Done)
	}
	t.Cleanup(func() {
		_, err := auth.NotionClient.Page.Update(context.Background(), notionapi.PageID(pageId), &notionapi.PageUpdateRequest{
			Archived:   true,
			Properties: notionapi.Properties{},
		})
		if err != nil {
			t.Error(err)
		}
	})
}

func TestGetTaskDetailsWithNullDate(t *testing.T) {
	details := test.CreateDetails()
	details.DueDate = nil
	pageId, _, err := Service.Insert(connection.NotionDatabaseId, &details)
	if err != nil {
		t.Error(err)
	}
	taskDetails, _, err := Service.GetTaskDetails(connection.NotionDatabaseId, pageId)
	if err != nil {
		t.Error(err)
	}
	if taskDetails.Title != details.Title {
		t.Errorf("Title is not correct: new value %v correct value %v", taskDetails.Title, details.Title)
	}
	if taskDetails.DueDate != nil {
		t.Errorf("Deadline is not correct: new value %v correct value %v", taskDetails.DueDate.Format(time.RFC3339), details.DueDate.Format(time.RFC3339))
	}
	if taskDetails.Done != details.Done {
		t.Errorf("Done is not correct: new value %v correct value %v", taskDetails.Done, details.Done)
	}
	t.Cleanup(func() {
		_, err := auth.NotionClient.Page.Update(context.Background(), notionapi.PageID(pageId), &notionapi.PageUpdateRequest{
			Archived:   true,
			Properties: notionapi.Properties{},
		})
		if err != nil {
			t.Error(err)
		}
	})
}

func TestCreateNotes(t *testing.T) {
	prop := notionapi.Properties{
		"Name1": NewRichTextProperty("Test"),
		"Name2": NewRichTextProperty("Test2"),
	}
	notes := createNotes(prop)
	notes2 := keys.BREAK_LINE + "Name1: Test\nName2: Test2\n"
	notes3 := keys.BREAK_LINE + "Name2: Test2\nName1: Test\n"
	if notes != notes2 && notes != notes3 {
		t.Errorf("Notes are not correct: new value \n %v correct value \n %v", notes, notes2)
	}
}

func TestUpdate(t *testing.T) {
	details := test.CreateDetails()
	pageId, _, err := Service.Insert(connection.NotionDatabaseId, &details)
	if err != nil {
		t.Error(err)
	}
	details.Title = "New title"
	newDate := details.DueDate.AddDate(0, 0, 1)
	details.DueDate = &newDate
	details.Done = true
	_, err = Service.Update(connection.NotionDatabaseId, pageId, &details)
	if err != nil {
		t.Error(err)
	}
	taskDetails, _, err := Service.GetTaskDetails(connection.NotionDatabaseId, pageId)
	if err != nil {
		t.Error(err)
	}
	if taskDetails.Title != details.Title {
		t.Errorf("Title is not correct: new value %v correct value %v", taskDetails.Title, details.Title)
	}
	if taskDetails.DueDate.Format(time.RFC3339) != details.DueDate.Format(time.RFC3339) {
		t.Errorf("Deadline is not correct: new value %v correct value %v", taskDetails.DueDate.Format(time.RFC3339), details.DueDate.Format(time.RFC3339))
	}
	if taskDetails.Done != details.Done {
		t.Errorf("Done is not correct: new value %v correct value %v", taskDetails.Done, details.Done)
	}
	t.Cleanup(func() {
		_, err := auth.NotionClient.Page.Update(context.Background(), notionapi.PageID(pageId), &notionapi.PageUpdateRequest{
			Archived:   true,
			Properties: notionapi.Properties{},
		})
		if err != nil {
			t.Error(err)
		}
	})
}
