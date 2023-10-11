package models

import (
	"time"
)

type ConnectedTask struct {
	TasksId      string `gorm:"primaryKey;uniqueIndex"`
	TaskUpdate   *time.Time
	NotionId     string `gorm:"primaryKey;uniqueIndex"`
	NotionUpdate *time.Time
	Connection   Connection `gorm:"foreignKey:NotionDatabasId;references:NotionDatabasId"`
	TasksListId  string
}

type Connection struct {
	NotionDatabasId string
	TasksListId     string
}
