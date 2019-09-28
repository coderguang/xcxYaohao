package httpHandle

import (
	"encoding/json"
	"net/http"
	"xcxYaohao/src/config"
	"xcxYaohao/src/data"

	"github.com/coderguang/GameEngine_go/sglog"
	"github.com/coderguang/GameEngine_go/sgtime"
	"github.com/coderguang/GameEngine_go/sgwx/sgwxopenid"
)

func logicHandle(w http.ResponseWriter, r *http.Request, flag chan bool) {

	returnData := make(map[string]interface{})
	returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_DO_NOT_THING

	defer func() {

		str, err := json.Marshal(returnData)
		if err != nil {
			sglog.Error("parse returnData to string error", err)
			return
		}
		sglog.Info("return str is", string(str))
		w.Write([]byte(string(str)))

		flag <- true
	}()

	r.ParseForm()

	sglog.Debug("require data is ", r.Form)

	op := r.FormValue(HTTP_ARGS_KEY)
	city := r.FormValue(HTTP_ARGS_CITY)
	loginCode := r.FormValue(HTTP_ARGS_CODE)

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
			returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_WX_ERROR_CODE
			return
		}
		data.AddWxOpenId(openId)

		_, err := data.GetNoticeData(openId.Openid)
		if err != nil {
			data.AddStatistic(data.StatisticNewOpenTimes, 1)
		}
	}

	//openId.Openid = loginCode

	if !config.IsSupportCity(city) && !isIgnoreCityArgs(op) {
		returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_TITLE
	} else {
		switch op {
		case HTTP_ARGS_TIME:
			// ?op=time&city=guangzhou&code=0
			requireLastestTime(r, city, openId.Openid, returnData)
		case HTTP_ARGS_SEARCH:
			// ?op=search&city=guangzhou&key=0000100748077&code=0
			matchData(r, city, openId.Openid, returnData)
		case HTTP_ARGS_BIND_GET_DATA:
			// ?op=getData&city=guangzhou&code=0
			getBindData(r, city, openId.Openid, returnData)
		case HTTP_ARGS_BIND_REQUIRE:
			// ?op=require&city=guangzhou&code=0&card=1&phone=
			requireRandomCode(r, city, openId.Openid, returnData)
		case HTTP_ARGS_BIND_CONFIRM:
			// ?op=confirm&&code=0&data=
			confirmRandomCode(r, city, openId.Openid, returnData)
		case HTTP_ARGS_BIND_CANCEL:
			// ?op=cancel&&code=0
			cancelBind(r, city, openId.Openid, returnData)
		case HTTP_ARGS_SHARE:
			// ?op=share&&code=
			share(r, city, openId.Openid, returnData)
		}
	}
}

func isIgnoreCityArgs(op string) bool {
	if op == HTTP_ARGS_BIND_CONFIRM || op == HTTP_ARGS_BIND_CANCEL || op == HTTP_ARGS_SHARE {
		return true
	}
	return false
}
