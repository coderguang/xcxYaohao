package httpHandle

import (
	"net/http"
	"xcxYaohao/src/data"
	"xcxYaohao/src/db"
	"xcxYaohao/src/define"
)

func requireLastestTime(r *http.Request, city string, openId string, returnData map[string]interface{}) {
	timestr := data.GetLastestInfo(city)
	returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_OK
	returnData[HTTP_RETURN_TIME] = timestr
	returnData[HTTP_RETURN_TIPS] = data.GetBoardcast(city)

	existData := data.AddOpenXcxTimes(openId, city)
	db.UpdateNoticeData(existData)

	data.AddStatistic(define.StatisticOpenTimes, 1)
}
