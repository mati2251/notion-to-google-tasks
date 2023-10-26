package db

import (
	"database/sql"
	"errors"
	"os"

	"github.com/mati2251/notion-to-google-tasks/keys"
	"github.com/mati2251/notion-to-google-tasks/models"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func OpenFile() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	filepath := home + keys.FILES_PATH + keys.DB_FILE
	_, errIfExist := os.Stat(filepath)
	DB, err = sql.Open("sqlite3", filepath)
	if err != nil {
		return err
	}
	if os.IsNotExist(errIfExist) {
		err = initDatabase()
		if err != nil {
			return err
		}
	}
	return nil
}

func initDatabase() error {
	_, err := DB.Exec("CREATE TABLE tasks (taskId varchar PRIMARY KEY UNIQUE, notionId varchar UNIQUE, taskUpdate datetime, notionUpdate datetime, connectionId varchar)")
	return err
}

func Insert(connectedTask models.ConnectedTask) error {
	_, err := DB.Exec(
		"INSERT INTO tasks(taskId, notionId, taskUpdate, notionUpdate, connectionId) VALUES (?, ?, ?, ?, ?)",
		connectedTask.TasksId,
		connectedTask.NotionId,
		connectedTask.TaskUpdate,
		connectedTask.NotionUpdate,
		connectedTask.ConnectionId,
	)
	if err != nil {
		return errors.Join(err, errors.New("error while inserting task"))
	}
	return nil
}