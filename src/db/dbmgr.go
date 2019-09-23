package db

import (
	"xcxYaohao/src/data"
	"xcxYaohao/src/define"

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
}

func initAndLoadDownloadHistory() {
	err := globalDb.AutoMigrate(define.DownloadHistoryUrl{}).Error
	if err != nil {
		sglog.Error("initAndLoadDownloadHistory", err)
	}

	sglog.Info("init and load hitory data ok")

	historyDatas := []define.DownloadHistoryUrl{}
	err = globalDb.Find(&historyDatas).Error
	if err != nil {
		sglog.Error("initAndLoadDownloadHistory find error", err)
	}
	data.InitDataFromDb(historyDatas)
}
