package models

import (
	"time"
)

type ConnectedTask struct {
	TasksId      string
	TaskUpdate   *time.Time
	NotionId     string
	NotionUpdate *time.Time
	ConnectionId string
}

type Connection struct {
	NotionDatabasId string
	TasksListId     string
}

type TaskDetails struct {
	Title   string
	DueDate *time.Time
	Done    bool
	Notes   string
}

type Service interface {
	Insert(connectionId string, details *TaskDetails) (string, *time.Time, error)
	Update(connectionId string, id string, details *TaskDetails) (*time.Time, error)
	GetTaskDetails(connectionId string, id string) (*TaskDetails, *time.Time, error)
}
