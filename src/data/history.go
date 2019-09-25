package data

import (
	"time"
	"xcxYaohao/src/define"

	"github.com/coderguang/GameEngine_go/sglog"
	"github.com/mohae/deepcopy"
)

func init() {
	globalHistoryData = new(define.SecureDownloadHistoryUrl)
	globalHistoryData.Data = make(map[string](map[string]*define.DownloadHistoryUrl))
}

var (
	globalHistoryData *define.SecureDownloadHistoryUrl
)

func InitHistoryDataFromDb(datas []define.DownloadHistoryUrl) {
	for _, v := range datas {
		tmp := deepcopy.Copy(v)
		tmpV, ok := tmp.(define.DownloadHistoryUrl)
		if !ok {
			sglog.Error("deep copy error,title:", v.Title, ",url:", v.URL)
			continue
		}
		if cityMap, ok := globalHistoryData.Data[v.Title]; ok {
			cityMap[v.URL] = &tmpV
		} else {
			globalHistoryData.Data[v.Title] = make(map[string]*define.DownloadHistoryUrl)
			globalHistoryData.Data[v.Title][v.URL] = &tmpV
		}
	}
	sglog.Info("load history from db complete,size=", len(datas))
}

func GetReDownloadList(title string) []*define.DownloadHistoryUrl {
	globalHistoryData.Lock.Lock()
	defer globalHistoryData.Lock.Unlock()
	relist := []*define.DownloadHistoryUrl{}
	if cityMap, ok := globalHistoryData.Data[title]; ok {
		for _, v := range cityMap {
			if v.Status != define.DEF_DOWNLOAD_STATUS_COMPLETE {
				relist = append(relist, v)
			}
		}
		return relist
	}
	return relist
}

func NeedDownloadFile(title string, url string) bool {
	globalHistoryData.Lock.Lock()
	defer globalHistoryData.Lock.Unlock()
	if cityMap, ok := globalHistoryData.Data[title]; ok {
		if v, ok := cityMap[url]; ok {
			if v.Status != define.DEF_DOWNLOAD_STATUS_COMPLETE {
				return true
			} else {
				return false
			}
		}
	}
	return true
}

func ChangeDownloadStatus(title string, url string, status int, tips string) *define.DownloadHistoryUrl {

	sglog.Info("change downstaus,", title, "===========>", url)

	globalHistoryData.Lock.Lock()
	defer globalHistoryData.Lock.Unlock()
	if cityMap, ok := globalHistoryData.Data[title]; ok {
		if v, ok := cityMap[url]; ok {
			v.Status = status
			v.Tips = tips
			return v
		} else {
			tmp := new(define.DownloadHistoryUrl)
			tmp.Title = title
			tmp.URL = url
			tmp.Status = status
			now := time.Now()
			tmp.DownloadDt = now
			tmp.Tips = tips
			globalHistoryData.Data[title][tmp.URL] = tmp
			return tmp
		}
	} else {
		globalHistoryData.Data[title] = make(map[string]*define.DownloadHistoryUrl)
		tmp := new(define.DownloadHistoryUrl)
		tmp.Title = title
		tmp.URL = url
		tmp.Status = status
		now := time.Now()
		tmp.DownloadDt = now
		tmp.Tips = tips
		globalHistoryData.Data[title][tmp.URL] = tmp
		return tmp
	}
}
