package httpHandle

import (
	"net/http"
	"xcxYaohao/src/data"

	"github.com/coderguang/GameEngine_go/sgtime"
	"github.com/coderguang/GameEngine_go/sgwx/sgwxopenid"
)

func requireLastestTime(r *http.Request, returnData map[string]interface{}) {
	loginCode := r.FormValue(HTTP_OP_CODE)

	//init wx openid
	openId, err := data.GetWxOpenId(loginCode)

	if err != nil {
		//new require
		openId = new(sgwxopenid.SWxOpenid)
		openId.Code = loginCode
		openId.Time = sgtime.New()
		appid, secret := data.GetAppidCfg()
		openId.Openid, err = sgwxopenid.GetOpenIdFromWx(appid, secret, openId.Code)
		if err != nil {
			returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_NOTICE_ERR_WX_ERROR_CODE
			return
		}
		data.AddWxOpenId(openId)
	}
}
