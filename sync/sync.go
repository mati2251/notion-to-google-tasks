package sync

import (
	"database/sql"
	"os"

	"github.com/mati2251/notion-to-google-tasks/keys"
	"github.com/mati2251/notion-to-google-tasks/models"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Sync(connections []models.Connection) error {
	err := openFile()
	if err != nil {
		return err
	}
	ids, err := updates()
	err = inserts(ids, "1")
	if err != nil {
		return err
	}

	return nil
}

func openFile() error {
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
