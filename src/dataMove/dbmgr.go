package dataMove

import (
	"strconv"
	"xcxYaohao/src/data"
	"xcxYaohao/src/db"

	"github.com/coderguang/GameEngine_go/sgcfg"
	"github.com/coderguang/GameEngine_go/sgdb/sgmysql"
	"github.com/coderguang/GameEngine_go/sglog"
	"github.com/coderguang/GameEngine_go/sgthread"
	"github.com/jinzhu/gorm"
)

var (
	globalDb *gorm.DB
)

func InitDb(cmd []string) {

	filename := sgcfg.GetServerCfgDir() + "oldmysql.json"
	cfg, err := sgmysql.ReadCfg(filename)
	if err != nil {
		sglog.Error("init old mysql db error,err:", err)
		sgthread.DelayExit(2)
	}

	globalDb, err = sgmysql.Open(cfg)
	if err != nil {
		sglog.Error("open old mysql db error,err:", err)
		sgthread.DelayExit(2)
	}

	sglog.Info("init old db connection ok")

	initOldRequireDatata()

	initOldNotice()
}

func initOldRequireDatata() {
	sglog.Info("star load old require data")

	datas := []XcxCardNoticeRequireData{}
	err := globalDb.Find(&datas).Error
	if err != nil {
		sglog.Error("init old notice find error", err)
	}
	sglog.Info("get notice data ok,size:", len(datas))
	for _, v := range datas {
		bindData, err := data.GetNoticeData(v.TokenId)
		if err != nil {
			bindData = data.AddOpenXcxTimes(v.TokenId)
			bindData.CreateDt = *v.FinalLogin
			bindData.RequireTimes = v.RequireTime
			bindData.ShareTimes = v.ShareTimes
			db.UpdateNoticeData(bindData)
		}
	}
}

func initOldNotice() {

	sglog.Info("star load old notice")

	datas := []XcxCardNotice{}
	err := globalDb.Find(&datas).Error
	if err != nil {
		sglog.Error("init old notice find error", err)
	}
	sglog.Info("get notice data ok,size:", len(datas))
	for _, v := range datas {
		bindData, err := data.GetNoticeData(v.TokenId)
		if err != nil {
			sglog.Info("transfrom notice record error,token:", v.TokenId)
			bindData = data.AddOpenXcxTimes(v.TokenId)
		}
		bindData.Title = v.Title
		bindData.CardType = v.CardType
		bindData.Code = v.Code
		bindData.Phone = v.Phone
		bindData.EndDt = *v.EndDt
		bindData.Desc = v.Tips
		bindData.RenewTimes = v.RenewTimes
		bindData.Status, _ = strconv.Atoi(v.Status)
		bindData.NoticeTimes = v.NoticeTimes
		db.UpdateNoticeData(bindData)
	}
}
