package data

import (
	"xcxYaohao/src/define"

	"github.com/coderguang/GameEngine_go/sglog"
)

func init() {
	globalHistoryData = new(define.SecureDownloadHistoryUrl)
}

var (
	globalHistoryData *define.SecureDownloadHistoryUrl
)

func InitDataFromDb(datas []define.DownloadHistoryUrl) {
	for _, v := range datas {
		if cityMap, ok := globalHistoryData.Data[v.Title]; ok {
			cityMap[v.URL] = &v
		} else {
			globalHistoryData.Data[v.Title] = make(map[string]*define.DownloadHistoryUrl)
			globalHistoryData.Data[v.Title][v.URL] = &v
		}
	}
	sglog.Info("load history from db complete,size=", len(datas))
}
