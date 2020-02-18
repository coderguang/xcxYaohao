package httpHandle

import (
	"net/http"
	"xcxYaohao/src/data"
	"xcxYaohao/src/db"
	"xcxYaohao/src/define"

	"github.com/coderguang/GameEngine_go/sglog"
)

func requireLastestTime(r *http.Request, city string, openId string, platform string, returnData map[string]interface{}) {
	timestr := data.GetLastestInfo(city)
	returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_OK
	returnData[HTTP_RETURN_TIME] = timestr
	returnData[HTTP_RETURN_TIPS] = data.GetBoardcast(city)
	if platform == "alipay" {
		returnData[HTTP_RETURN_TIPS] = ""
	}
	returnData[HTTP_RETURN_Data] = openId

	scenId := r.FormValue(HTTP_ARGS_SCENE_ID)
	shareBy := r.FormValue(HTTP_ARGS_SHARE_FROM)

	existData, shareByData := data.AddOpenXcxTimes(platform, openId, city, scenId, shareBy)
	db.UpdateNoticeData(existData)
	//sglog.Info("scenEId:", scenId, ",shareBy:", shareBy)
	if shareByData != nil && shareByData.Token != "" {
		shareByData.ShareToNum++
		db.UpdateNoticeData(shareByData)
		sglog.Info("new player shared by other:", shareByData.Token, ",num:", shareByData.ShareToNum, ",scene:", scenId)
	}

	data.AddStatistic(define.StatisticOpenTimes, 1)

}
