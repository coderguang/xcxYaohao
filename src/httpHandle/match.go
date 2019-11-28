package httpHandle

import (
	"net/http"
	"xcxYaohao/src/cache"
	"xcxYaohao/src/data"
	"xcxYaohao/src/define"
)

func matchData(r *http.Request, city string, openId string, returnData map[string]interface{}) {
	key := r.FormValue(HTTP_ARGS_MATCH_KEY)
	if ok, v := cache.GetMatchData(city, key); ok {
		tmplist := []*define.CardDataForClient{}
		for _, vv := range v {
			tmplist = append(tmplist, vv.CardDataToClient())
		}
		returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_OK
		returnData[HTTP_RETURN_Data] = tmplist
	} else {
		returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_NO_MATCH_DATA
	}

	data.AddStatistic(define.StatisticRequireTimes, 1)
}
