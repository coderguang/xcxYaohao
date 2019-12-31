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

	scenId := r.FormValue(HTTP_ARGS_SCENE_ID)
	shareBy := r.FormValue(HTTP_ARGS_SHARE_FROM)

	existData, shareByData := data.AddOpenXcxTimes(openId, city, scenId, shareBy)
	db.UpdateNoticeData(existData)
	if shareByData != nil && shareByData.Token != "" {
		shareByData.ShareToNum++
		db.UpdateNoticeData(shareByData)
	}

	data.AddStatistic(define.StatisticOpenTimes, 1)
}
