package sync

import (
	"time"

	"github.com/mati2251/notion-to-google-tasks/keys"

	"github.com/spf13/viper"
)

func SetLastTimeSync() {
	viper.Set(keys.LAST_TIME_SYNC, time.Now())
	viper.WriteConfig()
}

func GetLastTimeSync() time.Time {
	return viper.GetTime(keys.LAST_TIME_SYNC)
}
