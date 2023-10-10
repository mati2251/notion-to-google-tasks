package models

import (
	"time"
)

type ConnectedTask struct {
	TasksId      string
	TaskUpdate   time.Time
	NotionId     string
	NotionUpdate time.Time
	Connection   *Connection
}

type Connection struct {
	NotionDatabasId string
	TasksListId     string
}
