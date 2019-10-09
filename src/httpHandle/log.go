package httpHandle

import (
	"net/http"
	"xcxYaohao/src/data"
	"xcxYaohao/src/db"
	"xcxYaohao/src/define"
)

func share(r *http.Request, city string, openId string, returnData map[string]interface{}) {
	existData, err := data.GetNoticeData(openId)
	if err != nil {
		returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_NOT_BIND_DATA
		return
	}
	existData.ShareTimes++
	db.UpdateNoticeData(existData)

	data.AddStatistic(define.StatisticShareTime, 1)
	returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_OK
}
