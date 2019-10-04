package httpHandle

import (
	"time"
	"xcxYaohao/src/data"
	"xcxYaohao/src/db"
	"xcxYaohao/src/define"

	"github.com/coderguang/GameEngine_go/sglog"
)

func NoticeCurrentMonthDataUpdate(title string, curTime string) {
	finalNotice := data.GetNoticeFinalTime(title)
	if finalNotice >= curTime {
		sglog.Info("notice already ", title, curTime, ",last notice is ", finalNotice)
		return
	}
	finalNoticeData := data.UpdateNoticeFinalTime(title, curTime)
	err := db.UpdateNoticeFinalData(finalNoticeData)
	if err != nil {
		sglog.Error("update final notice time to db error,", err)
	}
	lucklist, unlucklist := data.GetSmsNoticeData(title)
	luckPhone := []string{}
	unluckPhone := []string{}

	now := time.Now()
	for _, v := range lucklist {
		bindData, err := data.GetNoticeData(v)
		if err != nil {
			continue
		}
		if bindData.Title != title {
			sglog.Error("no match title,need:", title, ",now:", bindData.Title, ",openid:", v)
			continue
		}
		if bindData.Status != define.YAOHAO_NOTICE_STATUS_NORMAL {
			continue
		}
		luckPhone = append(luckPhone, bindData.Phone)
		bindData.FinalNoticeDt = now
		bindData.Status = define.YAOHAO_NOTICE_STATUS_CANCEL_BY_GM_BECASURE_LUCK
		bindData.NoticeTimes++
		db.UpdateNoticeData(bindData)
	}

	sglog.Info("send sms luck to ", title, luckPhone)

	for _, v := range unlucklist {
		bindData, err := data.GetNoticeData(v)
		if err != nil {
			continue
		}
		if bindData.Title != title {
			sglog.Error("no match title,need:", title, ",now:", bindData.Title, ",openid:", v)
			continue
		}
		if bindData.Status != define.YAOHAO_NOTICE_STATUS_NORMAL {
			continue
		}
		unluckPhone = append(unluckPhone, bindData.Phone)
		bindData.FinalNoticeDt = now
		if bindData.EndDt.Month() == now.Month() {
			bindData.Status = define.YAOHAO_NOTICE_STATUS_TIME_OUT
		}
		bindData.NoticeTimes++
		db.UpdateNoticeData(bindData)
	}

	sglog.Info("send sms unluck to ", title, luckPhone)

}
