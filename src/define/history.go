package define

import (
	"sync"
	"time"
)

const (
	DEF_DOWNLOAD_STATUS_DOWNING  int = 0
	DEF_DOWNLOAD_STATUS_COMPLETE int = 1
	DEF_DOWNLOAD_STATUS_ERROR    int = 2
)

type DownloadHistoryUrl struct {
	Title      string `gorm:"primary_key;type:varchar(64)"`
	URL        string `gorm:"primary_key;type:varchar(200)"`
	Status     int    `gorm:"default:0;not null"`
	DownloadDt time.Time
	Tips       string `gorm:"type:varchar(255)"`
}

type SecureDownloadHistoryUrl struct {
	Data map[string](map[string]*DownloadHistoryUrl)
	Lock sync.RWMutex
}

func (data *SecureDownloadHistoryUrl) String() string {
	str := "\n\n=============history data start============="
	for k, v := range data.Data {
		for _, vv := range v {
			str += "\n " + k + "--->" + vv.URL
		}
	}
	str += "=============history data end============="
	return str
}
