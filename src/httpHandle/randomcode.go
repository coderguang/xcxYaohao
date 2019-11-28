package httpHandle

import (
	"strconv"
	"time"
	"xcxYaohao/src/cache"
	"xcxYaohao/src/data"
	"xcxYaohao/src/db"
	"xcxYaohao/src/define"
	"xcxYaohao/src/sms"

	"github.com/coderguang/GameEngine_go/sgstring"
	"github.com/coderguang/GameEngine_go/sgtime"

	"github.com/coderguang/GameEngine_go/sglog"
	"github.com/coderguang/GameEngine_go/sgregex"
)

func requireRandomCodeFromClient(title string, openid string, cardType string, code string, phone string, leftTime string) (string, YaoHaoNoticeError) {

	randomCode := ""

	cardTypeInt, err := strconv.Atoi(cardType)
	if err != nil {
		sglog.Debug("require card type invalid")
		return randomCode, YAOHAO_ERR_HTTP_REQ_CARD_TYPE
	}

	if !checkCardTypeValid(title, cardTypeInt) {
		sglog.Debug("require card type invalid")
		return randomCode, YAOHAO_ERR_HTTP_REQ_CARD_TYPE
	}

	leftTimeInt, err := strconv.Atoi(leftTime)
	if err != nil {
		return randomCode, YAOHAO_ERR_LEFT_TIME
	}

	if !CheckCodeValid(title, code) {
		sglog.Debug("require  code invalid")
		return randomCode, YAOHAO_ERR_CODE
	}

	if !sgregex.CNMobile(phone) {
		sglog.Debug("require phone invalid")
		return randomCode, YAOHAO_ERR_PHONE
	}

	if !data.CanBindPhone(phone) {
		sglog.Debug("require phone invalid")
		return randomCode, YAOHAO_ERR_PHONE_BIND_TOO_MANY
	}

	if ok, _ := cache.GetMatchData(title, code); ok {
		return randomCode, YAOHAO_ERR_CODE_HAD_LUCK
	}

	existData, err := data.GetNoticeData(openid)
	if err != nil {
		return randomCode, YAOHAO_ERR_OPEN_ID_PARAM_NUM
	}
	if existData.Status == define.YAOHAO_NOTICE_STATUS_GM_LIMIT {
		return randomCode, YAOHAO_ERR_GM_LIMIT
	}

	if existData.IsStillValid() {
		if !existData.IsDataChange(title, code, phone, cardTypeInt) {
			return randomCode, YAOHAO_ERR_CODE_STILL_VALID
		}
	}
	now := sgtime.New()
	oldData := data.GetRequireData(openid)
	if oldData != nil {
		lastRequireDt := sgtime.TransfromTimeToDateTime(oldData.RequireDt)
		distanceTime := sgtime.GetTotalSecond(now) - sgtime.GetTotalSecond(lastRequireDt)
		if !oldData.IsDataChange(title, code, phone, cardTypeInt) {
			//data no change
			if distanceTime <= int64(define.YAOHAO_NOTICE_REQUIRE_VALID_TIME) {
				if oldData.Status == int(define.YaoHaoNoticeRequireStatus_Answer_Complete) {
					return randomCode, YAOHAO_ERR_REQUIRE_HAD_CONFIRM
				}
				if oldData.AnswerTimes >= define.YAOHAO_NOTICE_CONFIRM_TIMES {
					return randomCode, YAOHAO_ERR_CONFIRM_MORE_TIMES
				} else {
					return randomCode, YAOHAO_ERR_REQUIRE_WAIT_ANSWER
				}
			} else if distanceTime <= int64(define.YAOHAO_NOTICE_REQUIRE_UNLOCK_TIME) {
				if oldData.RequireTimes >= define.YAOHAO_NOTICE_REQUIRE_MAX_TIMES {
					return randomCode, YAOHAO_ERR_REQUIRE_HAD_LOCK
				}
			} else {
				oldData.RequireTimes = 0 //reset
				oldData.AnswerTimes = 0
			}
		} else {
			//有数据变更
			if distanceTime <= int64(define.YAOHAO_NOTICE_SMS_TIME_LIMIT) {
				return randomCode, YAOHAO_ERR_REQUIRE_HAD_LOCK
			}
			if distanceTime <= int64(define.YAOHAO_NOTICE_REQUIRE_UNLOCK_TIME) {
				if oldData.RequireTimes >= define.YAOHAO_NOTICE_REQUIRE_MAX_TIMES {
					return randomCode, YAOHAO_ERR_REQUIRE_HAD_LOCK
				}
			} else {
				oldData.RequireTimes = 0
				oldData.AnswerTimes = 0
			}
		}

		tmpRandomCode := sgstring.RandNumStringRunes(define.YAOHAO_NOTICE_RANDOM_NUM_LENGTH)

		err = sms.SendRandomCode(phone, title, code, tmpRandomCode)

		if err != nil {
			return randomCode, YAOHAO_ERR_SMS_RESULT_PARSE_ERROR
		}

		oldData.RequireDt = time.Now()
		oldData.Status = int(define.YaoHaoNoticeRequireStatus_Wait_Answer)
		oldData.Token = openid
		oldData.Code = code
		oldData.CardType = cardTypeInt
		oldData.Phone = phone
		oldData.LeftTime = leftTimeInt
		oldData.RandomNum = tmpRandomCode
		oldData.RequireTimes++
		randomCode = tmpRandomCode

		//todo

	} else {
		//针对绑定后取消的限制

		if !data.CanGetRequire(openid) {
			sglog.Info("require too fast,limit it", title, openid)
			return randomCode, YAOHAO_ERR_SMS_CLIENT
		}

		tmpRandomCode := sgstring.RandNumStringRunes(define.YAOHAO_NOTICE_RANDOM_NUM_LENGTH)

		err = sms.SendRandomCode(phone, title, code, tmpRandomCode)

		if err != nil {
			return randomCode, YAOHAO_ERR_SMS_RESULT_PARSE_ERROR
		}

		newRequireData := new(define.SRequireData)
		newRequireData.Title = title
		newRequireData.Status = int(define.YaoHaoNoticeRequireStatus_Wait_Answer)
		newRequireData.AnswerTimes = 0
		newRequireData.RequireDt = time.Now()
		newRequireData.Token = openid
		newRequireData.CardType = cardTypeInt
		newRequireData.Code = code
		newRequireData.Phone = phone
		newRequireData.LeftTime = leftTimeInt
		newRequireData.RandomNum = tmpRandomCode
		newRequireData.RequireTimes = 0

		randomCode = newRequireData.RandomNum
		data.AddOrUpdateRequireData(newRequireData)
	}
	data.AddRequireTimeLimits(openid)
	data.AddStatistic(define.StatisticRandomCodeSend, 1)
	return randomCode, YAOHAO_OK
}

func checkCardTypeValid(title string, cardTypeInt int) bool {
	switch title {
	case define.CITY_GUANGZHOU:
		return define.CARD_TYPE_NORMAL == cardTypeInt || define.CARD_TYPE_NEW_ENGINE == cardTypeInt
	case define.CITY_SHENZHEN:
		return define.CARD_TYPE_NORMAL == cardTypeInt || define.CARD_TYPE_NEW_ENGINE == cardTypeInt
	case define.CITY_HANGZHOU:
		return define.CARD_TYPE_NORMAL == cardTypeInt || define.CARD_TYPE_NEW_ENGINE == cardTypeInt
	case define.CITY_TIANJIN:
		return define.CARD_TYPE_NORMAL == cardTypeInt || define.CARD_TYPE_NEW_ENGINE == cardTypeInt
	case define.CITY_HAINAN:
		return define.CARD_TYPE_NORMAL == cardTypeInt || define.CARD_TYPE_NEW_ENGINE == cardTypeInt
	}
	return false
}

func CheckCodeValid(title string, code string) bool {
	switch title {
	case define.CITY_GUANGZHOU:
		if !sgregex.AllNum(code) {
			return false
		}
		if len(code) != 13 {
			return false
		}
		return true
	case define.CITY_SHENZHEN:
		if !sgregex.AllNum(code) {
			return false
		}
		if len(code) != 13 {
			return false
		}
		return true
	case define.CITY_HANGZHOU:
		if !sgregex.AllNum(code) {
			return false
		}
		if len(code) != 13 {
			return false
		}
		return true
	case define.CITY_TIANJIN:
		if !sgregex.AllNum(code) {
			return false
		}
		if len(code) != 13 {
			return false
		}
		return true
	case define.CITY_HAINAN:
		if !sgregex.AllNum(code) {
			return false
		}
		if len(code) != 13 {
			return false
		}
		return true
	}
	return false
}

func confirmRandomCodeFromClient(token string, randomCode string) YaoHaoNoticeError {
	oldData := data.GetRequireData(token)

	if oldData == nil {
		sglog.Debug("no require ,title:,token:,randomCode:", token, randomCode)
		return YAOHAO_ERR_CONFIRM_NOT_REQUIRE
	}

	now := sgtime.New()
	distance := sgtime.GetTotalSecond(now) - sgtime.GetTotalSecond(sgtime.TransfromTimeToDateTime(oldData.RequireDt))
	if distance <= int64(define.YAOHAO_NOTICE_REQUIRE_VALID_TIME) {
		if oldData.Status == int(define.YaoHaoNoticeRequireStatus_Answer_Complete) {
			sglog.Debug("had answer already ,title:,token:,randomCode:", token, randomCode)
			return YAOHAO_ERR_REQUIRE_HAD_CONFIRM
		}
	} else {
		return YAOHAO_ERR_HTTP_RANDOM_CODE_TIME_OUT
	}
	if oldData.AnswerTimes >= define.YAOHAO_NOTICE_CONFIRM_TIMES {
		return YAOHAO_ERR_CONFIRM_MORE_TIMES
	}
	oldData.AnswerTimes++
	if oldData.RandomNum != randomCode {
		sglog.Debug("error randomcode ,title:,token:,randomCode:", token, randomCode)
		return YAOHAO_ERR_CONFIRM_RANDOMCODE
	}
	//验证通过

	existData, err := data.GetNoticeData(token)
	if err != nil {
		return YAOHAO_ERR_OPEN_ID_PARAM_NUM
	}
	if existData.Status == define.YAOHAO_NOTICE_STATUS_GM_LIMIT {
		return YAOHAO_ERR_GM_LIMIT
	}

	if existData.IsStillValid() {
		//仍在有效期
		data.DelPhoneBind(existData.Phone)
	}
	data.AddPhoneBind(oldData.Phone)

	//更改数据
	existData.Code = oldData.Code
	existData.Phone = oldData.Phone
	firstOfMonth := time.Date(now.Year(), now.Month(), 0, 0, 0, 0, 0, now.Location())
	existData.EndDt = firstOfMonth.AddDate(0, oldData.LeftTime, 0)
	existData.CardType = oldData.CardType
	existData.Desc = ""
	existData.Title = oldData.Title
	existData.RenewTimes++
	existData.Status = define.YAOHAO_NOTICE_STATUS_NORMAL

	data.RemoveRequireData(token)

	db.UpdateNoticeData(existData)

	data.AddStatistic(define.StatisticBindTimes, 1)

	return YAOHAO_OK
}
