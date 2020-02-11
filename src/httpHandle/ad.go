package httpHandle

import (
	"net/http"
	"time"
	"xcxYaohao/src/data"
	"xcxYaohao/src/db"
	"xcxYaohao/src/define"
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
