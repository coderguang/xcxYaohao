package httpHandle

import (
	"encoding/json"
	"net/http"
	"xcxYaohao/src/config"
	"xcxYaohao/src/data"
	"xcxYaohao/src/define"

	"github.com/coderguang/GameEngine_go/sgbytedance/sgttopenid"
	"github.com/coderguang/GameEngine_go/sglog"
	"github.com/coderguang/GameEngine_go/sgqq/sgqqopenid"
	"github.com/coderguang/GameEngine_go/sgtime"
	"github.com/coderguang/GameEngine_go/sgwx/sgwxopenid"
)

func logicHandle(w http.ResponseWriter, r *http.Request, flag chan bool) {

	startDt := sgtime.New()
	requireData := make(map[string][]string)
	returnData := make(map[string]interface{})
	returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_DO_NOT_THING

	r.ParseForm()

	requireData = r.Form
	//sglog.Debug("require data is ", r.Form)

	op := r.FormValue(HTTP_ARGS_KEY)
	city := r.FormValue(HTTP_ARGS_CITY)
	loginCode := r.FormValue(HTTP_ARGS_CODE)
	platform := r.FormValue(HTTP_ARGS_PLATFORM)
	searchName := r.FormValue(HTTP_ARGS_MATCH_KEY)

	if len(requireData) == 1 && op == "" && loginCode == "" {
		sglog.Debug("require from alipay")
		platform = define.PLATFORM_ALIPAY
		jsonMap := make(map[string]interface{})
		for k, _ := range requireData {
			//sglog.Debug("key is ", k)
			if err := json.Unmarshal([]byte(k), &jsonMap); err != nil {
				returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_WX_ERROR_CODE
				return
			}
			break
		}

		opStr, ok := jsonMap[HTTP_ARGS_KEY]
		if !ok {
			returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_WX_ERROR_CODE
			return
		}
		op, ok = opStr.(string)
		if !ok {
			returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_WX_ERROR_CODE
			return
		}

		cityStr, ok := jsonMap[HTTP_ARGS_CITY]
		if !ok {
			returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_WX_ERROR_CODE
			return
		}
		city, ok = cityStr.(string)
		if !ok {
			returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_WX_ERROR_CODE
			return
		}

		loginCodeStr, ok := jsonMap[HTTP_ARGS_CODE]
		if !ok {
			returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_WX_ERROR_CODE
			return
		}
		loginCode, ok = loginCodeStr.(string)
		if !ok {
			returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_WX_ERROR_CODE
			return
		}
		matchKeyStr, ok := jsonMap[HTTP_ARGS_MATCH_KEY]
		if ok {
			matchkeyValue, ok := matchKeyStr.(string)
			if ok {
				searchName = matchkeyValue
			}
		}
	}

	defer func() {
		endDt := sgtime.New()
		useTime := sgtime.GetTotalSecond(endDt) - sgtime.GetTotalSecond(endDt)
		if useTime > 0 {
			sglog.Debug("handle use ", useTime, "start:", startDt, "--------->", endDt)
		}
		// strRequire, err := json.Marshal(requireData)
		// if err == nil {
		// 	sglog.Debug("require data is", string(strRequire))
		// } else {
		// 	sglog.Error("parse requireData to string error", err)
		// }

		str, err := json.Marshal(returnData)
		if err == nil {
			//sglog.Info("return str is", string(str))
			w.Write([]byte(string(str)))
		} else {
			sglog.Error("parse returnData to string error", err)
		}

		sglog.Info("platform:", "["+platform+"]", " token:", "["+loginCode+"]", ",op:", "["+op+"]", ",return code:", returnData[HTTP_RETURN_ERR_CODE])

		flag <- true
	}()

	//init wx openid
	openId, err := data.GetWxOpenId(loginCode)
	if err != nil {
		//new require
		openId = new(sgwxopenid.SWxOpenid)
		openId.Code = loginCode
		openId.Time = sgtime.New()

		if platform == "" || platform == define.PLATFORM_WEIXIN {
			platform = define.PLATFORM_WEIXIN
			appid, secret := data.GetAppidCfg()
			openId.Openid, err = sgwxopenid.GetOpenIdFromWx(appid, secret, openId.Code)
			if err != nil {
				returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_WX_ERROR_CODE
				return
			}
		} else if platform == define.PLATFORM_BYTEDANCE {
			appid, secret := data.GetByteDanceCfg()
			openId.Openid, err = sgttopenid.GetOpenIdFromServer(appid, secret, openId.Code)
			if err != nil {
				returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_WX_ERROR_CODE
				return
			}
		} else if platform == define.PLATFORM_QQ {
			appid, secret := data.GetQQCfg()
			openId.Openid, err = sgqqopenid.GetOpenIdFromServer(appid, secret, openId.Code)
			if err != nil {
				returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_WX_ERROR_CODE
				return
			}
		} else if platform == define.PLATFORM_ALIPAY {
			//sglog.Debug("alipay no token,code", openId.Code)
			openId.Openid = openId.Code
			switch op {
			case HTTP_ARGS_TIME:
			case HTTP_ARGS_SEARCH:
			case HTTP_ARGS_SHARE:
				break
			default:
				returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_WX_ERROR_CODE
				return
			}
		} else {
			returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_WX_ERROR_CODE
			return
		}

		data.AddWxOpenId(openId)

		_, err := data.GetNoticeData(openId.Openid)
		if err != nil {
			if platform != define.PLATFORM_ALIPAY {
				data.AddStatistic(define.StatisticNewOpenTimes, 1)
			}
		}
	}
	loginCode = openId.Openid

	requireData["token"] = []string{openId.Openid}

	//openId.Openid = loginCode

	if !config.IsSupportCity(city) && !isIgnoreCityArgs(op) {
		returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_TITLE
	} else {
		switch op {
		case HTTP_ARGS_TIME:
			// ?op=time&city=guangzhou&code=0
			requireLastestTime(r, city, openId.Openid, platform, returnData)
		case HTTP_ARGS_SEARCH:
			// ?op=search&city=guangzhou&key=0000100748077&code=0
			matchDataByName(r, city, openId.Openid, searchName, returnData)
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
		case HTTP_ARGS_COMPLETE_AD:
			completeAd(r, city, openId.Openid, returnData)
		case HTTP_ARGS_ONE_KEY_RE_BIND:
			reBindOneKey(r, city, openId.Openid, returnData)
		case HTTP_ARGS_DO_AD_TASK:
			taskAdComplete(r, city, openId.Openid, returnData)
		}
	}
}

func isIgnoreCityArgs(op string) bool {
	if op == HTTP_ARGS_BIND_CONFIRM || op == HTTP_ARGS_BIND_CANCEL || op == HTTP_ARGS_SHARE || op == HTTP_ARGS_ONE_KEY_RE_BIND || op == HTTP_ARGS_DO_AD_TASK {
		return true
	}
	return false
}
