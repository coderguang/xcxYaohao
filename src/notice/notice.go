package notice

import (
	"strconv"
	"time"
	"xcxYaohao/src/cache"
	"xcxYaohao/src/config"
	"xcxYaohao/src/data"
	"xcxYaohao/src/db"
	"xcxYaohao/src/define"
	"xcxYaohao/src/sms"

	"github.com/coderguang/GameEngine_go/sgtime"

	"github.com/coderguang/GameEngine_go/sgmail"

	"github.com/coderguang/GameEngine_go/sglog"
)

func NoticeCurrentMonthDataUpdate(title string, curTime string) {
	finalNotice := data.GetNoticeFinalTime(title)
	if finalNotice >= curTime {
		sglog.Info("notice already ", title, curTime, ",last notice is ", finalNotice)
		return
	}
	sglog.Info("start to send sms,", title, curTime)
	//sgmail.SendMail("start send mails", []string{config.GetUtilCfg().Receiver}, title+" "+curTime)
	finalNoticeData := data.UpdateNoticeFinalTime(title, curTime)
	err := db.UpdateNoticeFinalData(finalNoticeData)
	if err != nil {
		sglog.Error("update final notice time to db error,", err)
	}
	lucklist, unlucklist := cache.GetSmsNoticeDataList(title)
	luckPhone := []string{}
	unluckPhone := []string{}

	failedSend := 0

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
		if sgtime.GetTotalDay(sgtime.TransfromTimeToDateTime(bindData.FinalNoticeDt)) == sgtime.GetTotalDay(sgtime.TransfromTimeToDateTime(now)) {
			sglog.Debug("title:", title, ",token:", bindData.Token, ",today had send,last send time:", bindData.FinalNoticeDt)
			continue
		}

		luckPhone = append(luckPhone, bindData.Phone)

		err = sms.SendLuck(bindData.Phone, title, curTime)
		if err != nil {
			failedSend++
			sglog.Error("send result luck error,phone:", bindData.Phone, ",token:", bindData.Token, ",err:", err)
		}

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

		err = sms.SendUnLuck(bindData.Phone, title, curTime)
		if err != nil {
			failedSend++
			sglog.Error("send result unluck error,phone:", bindData.Phone, ",token:", bindData.Token, ",err:", err)
		}

		unluckPhone = append(unluckPhone, bindData.Phone)
		bindData.FinalNoticeDt = now
		if bindData.EndDt.Month() == now.Month() {
			bindData.Status = define.YAOHAO_NOTICE_STATUS_TIME_OUT
		}
		bindData.NoticeTimes++
		db.UpdateNoticeData(bindData)
	}

	sglog.Info("send sms unluck to ", title, unluckPhone)

	mailInfo := "title:" + title + "\ntime" + curTime +
		"toalSend:" + strconv.Itoa(len(luckPhone)+len(unluckPhone)) + "\n" +
		"failed:" + strconv.Itoa(failedSend) + "\n" +
		"luck:" + strconv.Itoa(len(luckPhone)) + "\n" +
		"unluck:" + strconv.Itoa(len(unluckPhone)) + "\n"

	sgmail.SendMail("sms result:"+title+" ,"+curTime, []string{config.GetUtilCfg().Receiver}, mailInfo)

}

func NoticeSmsByCmd(cmd []string) {
	title := cmd[1]
	curTime := cmd[2]
	NoticeCurrentMonthDataUpdate(title, curTime)
}
