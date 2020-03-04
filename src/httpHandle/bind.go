package httpHandle

import (
	"net/http"
	"time"
	"xcxYaohao/src/data"
	"xcxYaohao/src/db"
	"xcxYaohao/src/define"

	"github.com/coderguang/GameEngine_go/sgtime"
)

func getBindData(r *http.Request, city string, openId string, returnData map[string]interface{}) {
	datas, err := data.GetNoticeData(openId)

	if err != nil {
		returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_NOT_BIND_DATA
		return
	}
	if datas.Status == define.YAOHAO_NOTICE_STATUS_NOT_BIND {
		returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_NOT_BIND_DATA
		return
	}

	returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_OK
	returnData[HTTP_RETURN_STATUS] = datas.Status
	returnData[HTTP_ARGS_CODE] = datas.Code
	returnData[HTTP_ARGS_BIND_PHONE] = datas.Phone
	returnData[HTTP_RETURN_TIME] = sgtime.YearString(sgtime.TransfromTimeToDateTime(datas.EndDt)) + sgtime.MonthString(sgtime.TransfromTimeToDateTime(datas.EndDt))
	returnData[HTTP_ARGS_CITY] = datas.Title
}

func requireRandomCode(r *http.Request, city string, openId string, returnData map[string]interface{}) {

	cardType := r.FormValue(HTTP_ARGS_BIND_CARD_TYPE)
	code := r.FormValue(HTTP_ARGS_BIND_CODE)
	phone := r.FormValue(HTTP_ARGS_BIND_PHONE)
	leftTime := r.FormValue(HTTP_ARGS_TIME)

	_, errcode := requireRandomCodeFromClient(city, openId, cardType, code, phone, leftTime)
	if errcode != YAOHAO_OK {
		returnData[HTTP_RETURN_ERR_CODE] = errcode
		return
	}
	returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_OK
	//returnData[HTTP_RETURN_Data] = randomCode

}

func confirmRandomCode(r *http.Request, city string, openId string, returnData map[string]interface{}) {
	randomCode := r.FormValue(HTTP_ARGS_DATA)
	errcode := confirmRandomCodeFromClient(openId, randomCode)
	returnData[HTTP_RETURN_ERR_CODE] = errcode
}

func cancelBind(r *http.Request, city string, openId string, returnData map[string]interface{}) {
	existData, err := data.GetNoticeData(openId)
	if err != nil {
		returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_NOT_BIND_DATA
		return
	}
	if existData.Status == define.YAOHAO_NOTICE_STATUS_NOT_BIND {
		returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_NOT_BIND_DATA
		return
	}
	if existData.Status != define.YAOHAO_NOTICE_STATUS_NORMAL {
		returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_STATUS_NOT_NORMAL
		return
	}
	data.DelPhoneBind(existData.Phone)
	existData.Status = define.YAOHAO_NOTICE_STATUS_CANCEL
	db.UpdateNoticeData(existData)

	data.AddStatistic(define.StatisticBindCancel, 1)

	returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_OK
}

func reBindOneKey(r *http.Request, city string, openId string, returnData map[string]interface{}) {
	existData, err := data.GetNoticeData(openId)
	if err != nil {
		returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_NOT_BIND_DATA
		return
	}
	if existData.Status == define.YAOHAO_NOTICE_STATUS_NOT_BIND {
		returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_NOT_BIND_DATA
		return
	}
	if existData.Status == define.YAOHAO_NOTICE_STATUS_NORMAL {
		returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_CODE_STILL_VALID
		return
	}
	if existData.Status == define.YAOHAO_NOTICE_STATUS_CANCEL_BY_GM_BECASURE_LUCK {
		returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_CODE_HAD_LUCK
		return
	}

	if existData.Status != define.YAOHAO_NOTICE_STATUS_CANCEL || existData.Status != define.YAOHAO_NOTICE_STATUS_TIME_OUT {
		returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_NOT_BIND_DATA
		return
	}
	if !data.CanBindPhone(existData.Phone) {
		returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_ERR_PHONE_BIND_TOO_MANY
		return
	}
	data.AddPhoneBind(existData.Phone)

	//更改数据
	now := time.Now()
	leftTimeInt := 3
	if existData.AdTimes > 0 && sgtime.GetTotalSecond(sgtime.TransfromTimeToDateTime(now))-sgtime.GetTotalSecond(sgtime.TransfromTimeToDateTime(existData.AdCompleteDt)) < int64(300) {
		leftTimeInt = 6
	}

	firstOfMonth := time.Date(now.Year(), now.Month(), 0, 0, 0, 0, 0, now.Location())
	existData.EndDt = firstOfMonth.AddDate(0, leftTimeInt, 0)
	existData.RenewTimes++
	existData.Status = define.YAOHAO_NOTICE_STATUS_NORMAL

	db.UpdateNoticeData(existData)

	data.AddStatistic(define.StatisticOneKeyReBind, 1)

	returnData[HTTP_RETURN_ERR_CODE] = YAOHAO_OK

}
