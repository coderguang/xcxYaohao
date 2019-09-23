package define

import (
	"sync"
	"time"
)

type DownloadHistoryUrl struct {
	Title      string `gorm:"primary_key;type:varchar(64)"`
	URL        string `gorm:"primary_key;type:varchar(900)"`
	Status     int    `gorm:"default:1;not null"`
	DownloadDt *time.Time
	Tips       string `gorm:"type:varchar(255)"`
}

type SecureDownloadHistoryUrl struct {
	Data map[string](map[string]*DownloadHistoryUrl)
	Lock sync.RWMutex
}
