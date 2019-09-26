package define

import (
	"sync"
	"time"

	"github.com/coderguang/GameEngine_go/sgtime"

	"github.com/coderguang/GameEngine_go/sglog"
)

const YAOHAO_NOTICE_SMS_TIME_LIMIT int = 60
const YAOHAO_NOTICE_SMS_TIME_LIMIT_30 int = 1800
const YAOHAO_NOTICE_REQUIRE_VALID_TIME int = 300
const YAOHAO_NOTICE_REQUIRE_UNLOCK_TIME int = 1800
const YAOHAO_NOTICE_CONFIRM_TIMES int = 3
const YAOHAO_NOTICE_RANDOM_NUM_LENGTH int = 4
const YAOHAO_NOTICE_REQUIRE_MAX_TIMES int = 3
const YAOHAO_NOTICE_PHONE_CAN_BIND_TOKEN_MAX int = 3

const (
	YAOHAO_NOTICE_STATUS_NORMAL                     int = 0
	YAOHAO_NOTICE_STATUS_CANCEL                     int = 1
	YAOHAO_NOTICE_STATUS_TIME_OUT                   int = 2
	YAOHAO_NOTICE_STATUS_GM_LIMIT                   int = 3
	YAOHAO_NOTICE_STATUS_CANCEL_BY_GM_BECASURE_LUCK int = 4
	YAOHAO_NOTICE_STATUS_NOT_BIND                   int = 5
)

type YaoHaoNoticeRequireStatus int

const (
	YaoHaoNoticeRequireStatus_Wait_Answer YaoHaoNoticeRequireStatus = iota
	YaoHaoNoticeRequireStatus_Answer_Complete
	YaoHaoNoticeRequireStatus_Wait_ReAnswer //应答错误再次等待
)

type NoticeData struct {
	Token        string `gorm:"primary_key;type:varchar(512)"`
	Status       int
	Name         string
	Title        string
	CardType     int
	Code         string
	Phone        string
	EndDt        time.Time
	Desc         string
	RenewTimes   int
	NoticeTimes  int
	RequireTimes int
	FinalLogin   time.Time
	CreateDt     time.Time
	ShareTimes   int
}

type SecureNoticeData struct {
	MapData map[string]*NoticeData
	Lock    sync.RWMutex
}

func (data *NoticeData) IsStillValid() bool {
	now := time.Now()
	if data.Status == YAOHAO_NOTICE_STATUS_NORMAL && now.Before(data.EndDt) {
		return true
	}
	return false
}

func (data *NoticeData) IsDataChange(title string, code string, phone string, cardType int) bool {
	if data.Title == title && data.Code == code && data.Phone == phone && data.CardType == cardType {
		return false
	}
	return true
}

type SRequireData struct {
	Token        string
	Title        string
	CardType     int
	Code         string
	Phone        string
	RequireDt    time.Time //请求时间
	RandomNum    string
	AnswerTimes  int //回应次数
	Status       int
	LeftTime     int
	RequireTimes int
}

type SecureSRequireData struct {
	Data map[string]*SRequireData
	Lock sync.RWMutex
}

func (data *SRequireData) IsDataChange(title string, code string, phone string, cardType int) bool {
	if data.Title == title && data.Code == code && data.Phone == phone && data.CardType == cardType {
		return false
	}
	return true
}

func (data *SRequireData) ShowMsg() {
	sglog.Debug("=======start===========")
	sglog.Debug("token:", data.Token)
	sglog.Debug("Title:", data.Title)
	sglog.Debug("Code:", data.Code)
	sglog.Debug("CardType:", data.CardType)
	sglog.Debug("Phone:", data.Phone)
	sglog.Debug("RequireDt:", data.RequireDt)
	sglog.Debug("RandomNum:", data.RandomNum)
	sglog.Debug("AnswerTimes:", data.AnswerTimes)
	sglog.Debug("status:", data.Status)
	sglog.Debug("requireTimes:", data.RequireTimes)
	sglog.Debug("=======end===========")
}

type SRequireLimit struct {
	RequireTimes  int
	RequireDt     *sgtime.DateTime
	LastRequireDt *sgtime.DateTime
}

type SecureRequireLimit struct {
	MapData map[string]*SRequireLimit
	Lock    sync.RWMutex
}

type SecurePhoneLimit struct {
	MapData map[string]int
	Lock    sync.RWMutex
}