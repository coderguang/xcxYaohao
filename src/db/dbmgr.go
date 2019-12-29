package db

import (
	"xcxYaohao/src/data"
	"xcxYaohao/src/define"

	"github.com/coderguang/GameEngine_go/sgtime"

	"github.com/coderguang/GameEngine_go/sgcfg"
	"github.com/coderguang/GameEngine_go/sgdb/sgmysql"
	"github.com/coderguang/GameEngine_go/sglog"
	"github.com/coderguang/GameEngine_go/sgthread"
	"github.com/jinzhu/gorm"
)

var (
	globalDb *gorm.DB
)

func InitDb() {

	cfg, err := sgmysql.ReadCfg(sgcfg.MySQLCfgFile)
	if err != nil {
		sglog.Error("init mysql db error,err:", err)
		sgthread.DelayExit(2)
	}

	globalDb, err = sgmysql.Open(cfg)
	if err != nil {
		sglog.Error("open mysql db error,err:", err)
		sgthread.DelayExit(2)
	}

	sglog.Info("init db connection ok")

	initAndLoadDownloadHistory()

	initAndLoadNoticeData()

	initAndLoadFinalNoticeTime()

}

func AutoMirgrateTable() {
	err := globalDb.AutoMigrate(define.NoticeFinalTime{}).Error
	if err != nil {
		sglog.Error("initAndLoadFinalNoticeTime", err)
	}

	err = globalDb.AutoMigrate(define.DownloadHistoryUrl{}).Error
	if err != nil {
		sglog.Error("initAndLoadDownloadHistory", err)
	}

	err = globalDb.AutoMigrate(define.CardData{}).Error
	if err != nil {
		sglog.Error("initAndLoadData", err)
	}

	err = globalDb.AutoMigrate(define.NoticeData{}).Error
	if err != nil {
		sglog.Error("initAndLoadNoticeData", err)
	}

	err = globalDb.AutoMigrate(define.StatisticsData{}).Error
	if err != nil {
		sglog.Error("init statis data error")
	}
}

func initAndLoadFinalNoticeTime() {

	sglog.Info("init and load initAndLoadFinalNoticeTime data ok")

	datas := []define.NoticeFinalTime{}
	err := globalDb.Find(&datas).Error
	if err != nil {
		sglog.Error("initAndLoadFinalNoticeTime find error", err)
	}
	for _, v := range datas {
		data.UpdateNoticeFinalTime(v.Title, v.Time)
	}
}

func initAndLoadDownloadHistory() {
	sglog.Info("init and load hitory data ok")

	historyDatas := []define.DownloadHistoryUrl{}
	err := globalDb.Find(&historyDatas).Error
	if err != nil {
		sglog.Error("initAndLoadDownloadHistory find error", err)
	}
	data.InitHistoryDataFromDb(historyDatas)
}

func InitAndLoadCardData() {
	//return
	sglog.Info("init and load initAndLoadCardData data ok")

	now := sgtime.New()
	cardDatas := []define.CardData{}
	err := globalDb.Find(&cardDatas).Error
	if err != nil {
		sglog.Error("initAndLoadCardData find error", err)
	}
	end := sgtime.New()
	useTime := sgtime.GetTotalSecond(end) - sgtime.GetTotalSecond(now)
	sglog.Info("load initAndLoadCardData data from db use ", useTime, " seconds,size:", len(cardDatas))

	now = sgtime.New()
	data.InitCardDataFromDb(cardDatas)
	end = sgtime.New()
	useTime = sgtime.GetTotalSecond(end) - sgtime.GetTotalSecond(now)

	sglog.Info("init initAndLoadCardData data to memory use ", useTime, " seconds,size:", len(cardDatas))
}

func initAndLoadNoticeData() {

	sglog.Info("init and load initAndLoadNoticeData data ok")

	now := sgtime.New()
	noticeDatas := []define.NoticeData{}
	err := globalDb.Find(&noticeDatas).Error
	if err != nil {
		sglog.Error("NoticeData find error", err)
	}
	end := sgtime.New()
	useTime := sgtime.GetTotalSecond(end) - sgtime.GetTotalSecond(now)
	sglog.Info("load NoticeData data from db use ", useTime, " seconds,size:", len(noticeDatas))

	now = sgtime.New()
	data.InitNoticeDataFromDb(noticeDatas)
	end = sgtime.New()
	useTime = sgtime.GetTotalSecond(end) - sgtime.GetTotalSecond(now)

	sglog.Info("init NoticeData data use ", useTime, " seconds,size:", len(noticeDatas))

}

func InitLastestCardData() {
	//-- select * from card_data group by time,title;
	lastsDatas := []define.CardData{}
	err := globalDb.Select(`*`).Table("card_data").Group("title,time,type,card_type").Scan(&lastsDatas).Error
	if err != nil {
		sglog.Error("init lastest card time data error,", err)
		sgthread.DelayExit(2)
		return
	}

	for _, v := range lastsDatas {
		data.UpdateLastestInfo(v.Title, v.CardType, v.Type, v.Time)
	}
	data.ShowLastestInfo([]string{})
}

func IsCardDataExistByDb(title string, code string) bool {
	cards := []define.CardData{}
	err := globalDb.Where("title =? AND code=?", title, code).Find(&cards).Error
	if err != nil {
		sglog.Error("find data exist error,", title, code, err)
		return false
	}
	if len(cards) == 0 {
		return false
	}
	return true
}

func GetMatchDataByDb(title string, code string) (bool, []*define.CardData) {
	cards := []*define.CardData{}
	err := globalDb.Where("title =? AND (code=? OR name=?)", title, code, code).Find(&cards).Error
	if err != nil {
		sglog.Error("find data exist error,", title, code, err)
		return false, cards
	}
	if len(cards) != 0 {
		return true, cards
	}
	return false, cards
}
