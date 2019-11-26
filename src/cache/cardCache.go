package cache

import (
	"xcxYaohao/src/config"
	"xcxYaohao/src/data"
	"xcxYaohao/src/db"
	"xcxYaohao/src/define"

	"github.com/coderguang/GameEngine_go/sglog"
)

func InitAndLoadCardData() {

	if config.IsNoCache() {
		sglog.Debug("no cache flag!true")
		db.InitLastestCardData()
	} else {
		db.InitAndLoadCardData()
		data.ShowLastestInfo([]string{})
	}
}

func IsCardDataExist(title string, code string) bool {
	if config.IsNoCache() {
		return db.IsCardDataExistByDb(title, code)
	} else {
		return data.IsDataExist(title, code)
	}
}

func AddCardDataToMem(datas map[string]*define.CardData) {
	if config.IsNoCache() {
		sglog.Debug("no cahce,data would not cache,size:", len(datas))
		return
	} else {
		data.AddCardData(datas)
	}
}

func GetMatchData(title string, code string) (bool, []*define.CardData) {
	if config.IsNoCache() {
		return db.GetMatchDataByDb(title, code)
	} else {
		return data.GetMatchData(title, code)
	}
}

func GetSmsNoticeDataList(title string) ([]string, []string) {
	if config.IsNoCache() {
		luckList := []string{}
		unluckList := []string{}

		checkDatas := data.GetNoticeOpenIdAndCodeMap(title)

		for k, v := range checkDatas {
			if db.IsCardDataExistByDb(title, v) {
				luckList = append(luckList, k)
			} else {
				unluckList = append(unluckList, k)
			}
		}
		return luckList, unluckList
	} else {
		return data.GetSmsNoticeData(title)
	}

}
