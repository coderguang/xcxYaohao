package httpHandle

import (
	"net/http"
	"xcxYaohao/src/cache"
	"xcxYaohao/src/data"
	"xcxYaohao/src/db"
	"xcxYaohao/src/define"

	"github.com/coderguang/GameEngine_go/sglog"
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

func matchDataByName(r *http.Request, city string, openId string, key string, returnData map[string]interface{}) {
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

func TestMathDataByDbSort(cmd []string) {
	ok, results := db.GetMatchDataByDb("guangzhou", "李静")
	if ok {
		for _, v := range results {
			sglog.Debug("code:", v.Code, ",name:", v.Name, ",time:", v.Time)
		}
	}

}
