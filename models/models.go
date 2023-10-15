package models

import (
	"time"
)

type ConnectedTask struct {
	TasksId      string
	TaskUpdate   *time.Time
	NotionId     string
	NotionUpdate *time.Time
	Connection   Connection
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
	Inserts(ids []string, connectionId string) error
}
