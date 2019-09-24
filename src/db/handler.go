package db

import (
	"xcxYaohao/src/define"

	"github.com/coderguang/GameEngine_go/sglog"
)

func UpdateDownloadToDb(data *define.DownloadHistoryUrl) {
	err := globalDb.Where(define.DownloadHistoryUrl{Title: data.Title, URL: data.URL}).Assign(define.DownloadHistoryUrl{Status: data.Status, Tips: data.Tips}).FirstOrCreate(data).Error
	if err != nil {
		sglog.Error("update download url db error,", err)
	}
}
