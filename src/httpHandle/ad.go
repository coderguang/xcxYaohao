package httpHandle

import (
	"math/rand"
	"net/http"
	"time"
	"xcxYaohao/src/data"
	"xcxYaohao/src/db"
	"xcxYaohao/src/define"

	"github.com/coderguang/GameEngine_go/sgtime"
)

func completeAd(r *http.Request, city string, openId string, returnData map[string]interface{}) {

	returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_OK

	datas, err := data.GetNoticeData(openId)

	if err != nil {
		return
	}
	datas.AdCompleteDt = time.Now()
	datas.AdTimes++

	db.UpdateNoticeData(datas)
	data.AddStatistic(define.StatisticAdComplete, 1)
}

func taskAdComplete(r *http.Request, city string, openId string, returnData map[string]interface{}) {
	existData, err := data.GetNoticeData(openId)
	if err != nil {
		return
	}
	now := time.Now()
	if existData.AdTaskTimes > 0 {
		nowTotalSecond := sgtime.GetTotalSecond(sgtime.TransfromTimeToDateTime(now))
		if nowTotalSecond-sgtime.GetTotalSecond(sgtime.TransfromTimeToDateTime(existData.AdTaskDt)) > define.YAOHAO_AD_TASK_DISTANCE {
			existData.AdTaskValidTimes++
			existData.AdTaskDt = now
			//增加额外时长
			if existData.Status == define.YAOHAO_NOTICE_STATUS_NORMAL {
				returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_DO_AD_TASK_OK
				endTotalSecond := sgtime.GetTotalSecond(sgtime.TransfromTimeToDateTime(existData.EndDt))
				if endTotalSecond-nowTotalSecond < define.YAOHAO_NOTICE_MAX_NORMAL_TIME {
					//少于12个月时，每次必定增加1个月
					existData.EndDt = existData.EndDt.AddDate(0, 1, 0)
				} else if endTotalSecond-nowTotalSecond < define.YAOHAO_NOTICE_MAX_LUCKK_TIME {
					returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_DO_AD_TASK_OK_AND_UNLUCK
					//大于12个月时，每次5%的概率直接变为2年
					if existData.AdTaskValidTimes > int(define.YAOHAO_NOTICE_LUCK_AD_TIMES_LIMIT) {
						rand.Seed(time.Now().Unix())
						luckNum := rand.Intn(define.YAOHAO_NOTICE_LUCK_BASE)
						if luckNum < define.YAOHAO_NOTICE_LUCK_RATE {
							firstOfMonth := time.Date(now.Year(), now.Month(), 0, 0, 0, 0, 0, now.Location())
							existData.EndDt = firstOfMonth.AddDate(2, 0, 0)
							returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_DO_AD_TASK_OK_AND_LUCK
						}
					}
				} else {
					returnData[HTTP_RETURN_Data] = YAOHAO_ERR_DO_AD_TASK_OK_BUG_MAX
				}
			} else {
				returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_NOT_BIND_DATA
			}
		} else {
			returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_DO_AD_TASK_TOO_SHORT
		}
	} else {
		existData.AdTaskDt = now
		existData.AdTaskValidTimes++
		if existData.Status == define.YAOHAO_NOTICE_STATUS_NORMAL {
			returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_DO_AD_TASK_OK
			existData.EndDt = existData.EndDt.AddDate(0, 1, 0)
		} else {
			returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_NOT_BIND_DATA
		}
	}
	existData.AdTaskTimes++
	db.UpdateNoticeData(existData)
	data.AddStatistic(define.StatisticDoAdTask, 1)
}
