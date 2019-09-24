package data

import (
	"xcxYaohao/src/define"

	"github.com/coderguang/GameEngine_go/sglog"
)

func init() {
	globalHistoryData = new(define.SecureDownloadHistoryUrl)
	globalIgnonreData = new(define.SecureIgnoreUrl)
	globalHadVisitData = new(define.SecureHadVisitUrl)
}

var (
	globalHistoryData  *define.SecureDownloadHistoryUrl
	globalIgnonreData  *define.SecureIgnoreUrl
	globalHadVisitData *define.SecureHadVisitUrl
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

func AddIgnoreURL(title, url string) {
	globalIgnonreData.Lock.Lock()
	defer globalIgnonreData.Lock.Unlock()

	if cityMap, ok := globalIgnonreData.Data[title]; ok {
		cityMap[url] = ""
	} else {
		globalIgnonreData.Data[title] = make(map[string]string)
		globalIgnonreData.Data[title][url] = ""
	}
}