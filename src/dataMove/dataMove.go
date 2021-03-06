package dataMove

import (
	"github.com/coderguang/GameEngine_go/sglog"
	"github.com/coderguang/GameEngine_go/sgtime"
)

type XcxCardNotice struct {
	TokenId     string `gorm:"primary_key"`
	Name        string
	Title       string
	CardType    int
	Code        string
	Phone       string
	EndDt       *sgtime.DateTime
	Tips        string
	RenewTimes  int
	Status      string
	NoticeTimes int
}

func (u XcxCardNotice) TableName() string {
	return "xcx_card_notice"
}

func (data *XcxCardNotice) ShowMsg() {
	sglog.Debug("token:", data.TokenId)
	sglog.Debug("Name:", data.Name)
	sglog.Debug("Title:", data.Title)
	sglog.Debug("CardType:", data.CardType)
	sglog.Debug("Code:", data.Code)
	sglog.Debug("Phone:", data.Phone)
	sglog.Debug("EndDt:", sgtime.NormalString(data.EndDt))
	sglog.Debug("Desc:", data.Tips)
	sglog.Debug("RenewTimes:", data.RenewTimes)
}

type XcxCardNoticeRequireData struct {
	TokenId     string `gorm:"primary_key"`
	Title       string
	RequireTime int
	FinalLogin  *sgtime.DateTime
	Name        string
	ShareTimes  int
	Tips        string
}

func (u XcxCardNoticeRequireData) TableName() string {
	return "xcx_card_notice_require_data"
}

type XcxCardDataGuangzhou struct {
	Type     int
	CardType int
	Code     string `gorm:"primary_key"`
	Name     string `gorm:"primary_key"`
	Time     string
	Tips     string
}

func (u XcxCardDataGuangzhou) TableName() string {
	return "xcx_card_data_guangzhou"
}

type XcxCardDataShenzhen struct {
	Type     int
	CardType int
	Code     string `gorm:"primary_key"`
	Name     string `gorm:"primary_key"`
	Time     string
	Tips     string
}

func (u XcxCardDataShenzhen) TableName() string {
	return "xcx_card_data_shenzhen"
}
