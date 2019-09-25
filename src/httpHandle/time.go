package httpHandle

import (
	"net/http"
	"xcxYaohao/src/data"
)

func requireLastestTime(r *http.Request, city string, openId string, returnData map[string]interface{}) {
	timestr := data.GetLastestInfo(city)
	returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_OK
	returnData[HTTP_RETURN_TIME] = timestr
}
