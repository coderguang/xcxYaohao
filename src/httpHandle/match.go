package httpHandle

import (
	"encoding/json"
	"net/http"
	"xcxYaohao/src/data"
	"xcxYaohao/src/define"
)

func matchData(r *http.Request, city string, openId string, returnData map[string]interface{}) {
	key := r.FormValue(HTTP_ARGS_MATCH_KEY)
	if ok, v := data.GetMatchData(city, key); ok {
		tmplist := []*define.CardDataForClient{}
		for _, vv := range v {
			tmplist = append(tmplist, vv.CardDataToClient())
		}
		jsonBytes, _ := json.Marshal(tmplist)
		returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_OK
		returnData[HTTP_RETURN_MATCH_Data] = string(jsonBytes)
	} else {
		returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_NO_MATCH_DATA
	}
}
