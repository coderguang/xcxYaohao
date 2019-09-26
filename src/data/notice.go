package data

import (
	"errors"
	"xcxYaohao/src/define"

	"github.com/coderguang/GameEngine_go/sgtime"

	"github.com/coderguang/GameEngine_go/sglog"
	"github.com/mohae/deepcopy"
)

var (
	globalNoticeData   *define.SecureNoticeData
	globalRequireData  *define.SecureSRequireData
	globalRequireLimit *define.SecureRequireLimit
	globalPhoneLimit   *define.SecurePhoneLimit
)

func init() {
	globalNoticeData = new(define.SecureNoticeData)
	globalRequireData = new(define.SecureSRequireData)
	globalRequireLimit = new(define.SecureRequireLimit)
	globalPhoneLimit = new(define.SecurePhoneLimit)

	globalNoticeData.MapData = make(map[string]*define.NoticeData)
	globalRequireData.Data = make(map[string]*define.SRequireData)
	globalRequireLimit.MapData = make(map[string]*define.SRequireLimit)
	globalPhoneLimit.MapData = make(map[string]int)
}

func InitNoticeDataFromDb(datas []define.NoticeData) {
	globalNoticeData.Lock.Lock()
	defer globalNoticeData.Lock.Unlock()

	globalPhoneLimit.Lock.Lock()
	defer globalPhoneLimit.Lock.Unlock()

	for _, v := range datas {
		tmp := deepcopy.Copy(v)
		tmpV, ok := tmp.(define.NoticeData)
		if !ok {
			sglog.Error("deepcopy notice error,token:", v.Token)
			continue
		}
		globalNoticeData.MapData[v.Token] = &tmpV
		if v.Status == define.YAOHAO_NOTICE_STATUS_NORMAL {
			globalPhoneLimit.MapData[v.Phone]++
		}

	}
}

func GetNoticeData(openid string) (*define.NoticeData, error) {
	globalNoticeData.Lock.Lock()
	defer globalNoticeData.Lock.Unlock()
	if v, ok := globalNoticeData.MapData[openid]; ok {
		return v, nil
	}
	return nil, errors.New("no bind data")
}

func CanBindPhone(phone string) bool {
	globalPhoneLimit.Lock.Lock()
	defer globalPhoneLimit.Lock.Unlock()
	if v, ok := globalPhoneLimit.MapData[phone]; ok {
		if v >= define.YAOHAO_NOTICE_PHONE_CAN_BIND_TOKEN_MAX {
			return false
		}
	}
	return true
}

func GetRequireData(openid string) *define.SRequireData {
	globalRequireData.Lock.Lock()
	defer globalRequireData.Lock.Unlock()
	if v, ok := globalRequireData.Data[openid]; ok {
		return v
	}
	return nil
}

func CanGetRequire(openid string) bool {
	globalRequireLimit.Lock.Lock()
	defer globalRequireLimit.Lock.Unlock()
	if v, ok := globalRequireLimit.MapData[openid]; ok {
		now := sgtime.New()
		distance := sgtime.GetTotalSecond(now) - sgtime.GetTotalSecond(sgtime.TransfromTimeToDateTime(*v.LastRequireDt))
		if distance <= int64(define.YAOHAO_NOTICE_SMS_TIME_LIMIT) {
			return false
		} else if distance <= int64(define.YAOHAO_NOTICE_SMS_TIME_LIMIT_30) {
			if v.RequireTimes >= define.YAOHAO_NOTICE_REQUIRE_MAX_TIMES {
				return false
			}
		}
	}
	return true
}

func AddOrUpdateRequireData(data *define.SRequireData) {
	globalRequireData.Lock.Lock()
	defer globalRequireData.Lock.Unlock()
	globalRequireData.Data[data.Token] = data
	sglog.Info("add or update require data complete")
	data.ShowMsg()
}

func AddRequireTimeLimits(openid string) {
	globalRequireLimit.Lock.Lock()
	defer globalRequireLimit.Lock.Unlock()
	now := sgtime.New()
	if v, ok := globalRequireLimit.MapData[openid]; ok {
		now := sgtime.New()
		distance := sgtime.GetTotalSecond(now) - sgtime.GetTotalSecond(sgtime.TransfromTimeToDateTime(*v.LastRequireDt))
		if distance >= int64(define.YAOHAO_NOTICE_SMS_TIME_LIMIT_30) {
			v.RequireTimes = 1
			v.RequireDt = now
		} else {
			v.RequireTimes++
		}
		v.LastRequireDt = now
	} else {
		newData := new(define.SRequireLimit)
		newData.RequireTimes = 1
		newData.RequireDt = now
		newData.LastRequireDt = now
		globalRequireLimit.MapData[openid] = newData
	}
}
