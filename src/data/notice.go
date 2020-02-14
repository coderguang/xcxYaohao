package data

import (
	"errors"
	"time"
	"xcxYaohao/src/define"

	"github.com/coderguang/GameEngine_go/sgtime"

	"github.com/coderguang/GameEngine_go/sglog"
	"github.com/mohae/deepcopy"
)

var (
	globalNoticeData      *define.SecureNoticeData
	globalRequireData     *define.SecureSRequireData
	globalRequireLimit    *define.SecureRequireLimit
	globalPhoneLimit      *define.SecurePhoneLimit
	globalFinalNoticeData *define.SecureNoticeFinalTime
)

func init() {
	globalNoticeData = new(define.SecureNoticeData)
	globalRequireData = new(define.SecureSRequireData)
	globalRequireLimit = new(define.SecureRequireLimit)
	globalPhoneLimit = new(define.SecurePhoneLimit)
	globalFinalNoticeData = new(define.SecureNoticeFinalTime)

	globalNoticeData.MapData = make(map[string]*define.NoticeData)
	globalRequireData.Data = make(map[string]*define.SRequireData)
	globalRequireLimit.MapData = make(map[string]*define.SRequireLimit)
	globalPhoneLimit.MapData = make(map[string]int)
	globalFinalNoticeData.Data = make(map[string]*define.NoticeFinalTime)
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
	return nil, errors.New("no this data")
}

func AddOpenXcxTimes(platform string, openid string, title string, scenId string, shareFrom string) (*define.NoticeData, *define.NoticeData) {
	data, err := GetNoticeData(openid)
	shareFromData := new(define.NoticeData)
	shareFromData.Token = ""
	if err != nil {
		data = new(define.NoticeData)
		data.Token = openid
		data.CreateDt = time.Now()
		data.RequireTimes = 1
		data.Status = define.YAOHAO_NOTICE_STATUS_NOT_BIND
		data.FinalNoticeDt = time.Now()
		data.Title = title
		data.SceneId = scenId
		data.SharedBy = shareFrom
		data.Platform = platform

		if "" != shareFrom && shareFrom != data.Token {
			shareFromData, err = GetNoticeData(shareFrom)
			if err != nil {
				sglog.Debug("share token not find,token:", "["+shareFrom+"]", ",err:", err)
			}
		}
		globalNoticeData.Lock.Lock()
		defer globalNoticeData.Lock.Unlock()

		globalNoticeData.MapData[openid] = data

	} else {
		data.RequireTimes++
	}
	data.FinalLoginDt = time.Now()
	if "" == data.Phone {
		data.Title = title
	}

	return data, shareFromData
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

func RemoveRequireData(openid string) {
	globalRequireData.Lock.Lock()
	defer globalRequireData.Lock.Unlock()
	if _, ok := globalRequireData.Data[openid]; ok {
		delete(globalRequireData.Data, openid)
	}
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

func AddPhoneBind(phone string) {
	globalPhoneLimit.Lock.Lock()
	defer globalPhoneLimit.Lock.Unlock()
	if v, ok := globalPhoneLimit.MapData[phone]; ok {
		v++
	} else {
		globalPhoneLimit.MapData[phone] = 1
	}
}

func DelPhoneBind(phone string) {
	globalPhoneLimit.Lock.Lock()
	defer globalPhoneLimit.Lock.Unlock()
	if v, ok := globalPhoneLimit.MapData[phone]; ok {
		v--
		if 0 == v {
			delete(globalPhoneLimit.MapData, phone)
		}
	} else {
		sglog.Error("try to delete unbind phone,phon:", phone)
	}
}

func UpdateNoticeFinalTime(title string, time string) *define.NoticeFinalTime {
	globalFinalNoticeData.Lock.Lock()
	globalFinalNoticeData.Lock.Unlock()
	if v, ok := globalFinalNoticeData.Data[title]; ok {
		v.Time = time
		return v
	} else {
		tmp := new(define.NoticeFinalTime)
		tmp.Title = title
		tmp.Time = time
		globalFinalNoticeData.Data[title] = tmp
		return tmp
	}
}

func GetNoticeFinalTime(title string) string {
	globalFinalNoticeData.Lock.Lock()
	globalFinalNoticeData.Lock.Unlock()
	if v, ok := globalFinalNoticeData.Data[title]; ok {
		return v.Time
	}
	return "201909"
}

func GetSmsNoticeData(title string) ([]string, []string) {
	globalNoticeData.Lock.Lock()
	defer globalNoticeData.Lock.Unlock()
	globalCardData.Lock.Lock()
	defer globalCardData.Lock.Unlock()
	luckList := []string{}
	unluckList := []string{}

	if cityCards, ok := globalCardData.Data[title]; ok {
		for k, v := range globalNoticeData.MapData {
			if v.Title != title {
				continue
			}
			if _, luckFlag := cityCards[v.Code]; luckFlag {
				luckList = append(luckList, k)
			} else {
				unluckList = append(unluckList, k)
			}
		}
	}

	return luckList, unluckList
}

func GetNoticeOpenIdAndCodeMap(title string) map[string]string {
	globalNoticeData.Lock.Lock()
	defer globalNoticeData.Lock.Unlock()

	datas := make(map[string]string)

	for k, v := range globalNoticeData.MapData {
		if v.Title != title {
			continue
		}
		if v.Status != define.YAOHAO_NOTICE_STATUS_NORMAL {
			continue
		}
		datas[k] = v.Code
	}
	return datas
}
