package sync

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/mati2251/notion-to-google-tasks/keys"
	"github.com/mati2251/notion-to-google-tasks/models"
)

var synchronizedTasksFile io.Reader

func Sync(connections []models.Connection) error {
	synchronizedTasksFile, err = csv.NewReader(synchronizedTasksFile).ReadAll()
	fmt.Println(test)
	return nil
}

func openFile() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	filepath := home + keys.FILES_PATH + keys.SYNCHRONIZED_TASK_FILE_NAME
	synchronizedTasksFile, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	return nil
}
