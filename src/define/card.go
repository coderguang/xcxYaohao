package define

import (
	"sync"
	"time"
)

const (
	CARD_TYPE_NORMAL      = 1
	CARD_TYPE_NEW_ENGINE  = 2
	MEMBER_TYPE_PERSIONAL = 1
	MEMBER_TYPE_COMPANY   = 2
)

const (
	CITY_GUANGZHOU string = "guangzhou"
	CITY_SHENZHEN  string = "shenzhen"
	CITY_HANGZHOU  string = "hangzhou"
	CITY_TIANJIN   string = "tianjin"
)

type SLastestCardData struct {
	Title                 string
	TimeStr               string
	PersonalNormalUpdate  bool
	PersonalJieNengUpdate bool
	CompanyNormalUpdate   bool
	CompanyJieNengUpdate  bool
}

func (data *SLastestCardData) Reset() {
	data.PersonalNormalUpdate = false
	data.PersonalJieNengUpdate = false
	data.CompanyNormalUpdate = false
	data.CompanyJieNengUpdate = false
}
func (data *SLastestCardData) UpdateLastestInfo(cardType int, memberType int) {
	if MEMBER_TYPE_PERSIONAL == memberType {
		if CARD_TYPE_NORMAL == cardType {
			data.PersonalNormalUpdate = true
		} else {
			data.PersonalJieNengUpdate = true
		}
	} else {
		if CARD_TYPE_NORMAL == cardType {
			data.CompanyNormalUpdate = true
		} else {
			data.CompanyJieNengUpdate = true
		}
	}
}

func (data *SLastestCardData) IsAllCardInfoUpdate() bool {
	switch data.Title {
	case CITY_GUANGZHOU:
		if data.PersonalJieNengUpdate && data.PersonalNormalUpdate && data.CompanyJieNengUpdate && data.CompanyNormalUpdate {
			return true
		}
	case CITY_SHENZHEN:
		if data.PersonalNormalUpdate && data.CompanyNormalUpdate {
			return true
		}
	case CITY_HANGZHOU:
		if data.PersonalNormalUpdate && data.CompanyNormalUpdate {
			return true
		}
	case CITY_TIANJIN:
		if data.PersonalJieNengUpdate && data.PersonalNormalUpdate && data.CompanyJieNengUpdate && data.CompanyNormalUpdate {
			return true
		}
	}
	return false
}

type SecureLastestCardData struct {
	Data map[string]*SLastestCardData
	Lock sync.RWMutex
}

type CardData struct {
	Title    string `gorm:"primary_key;type:varchar(64)"`
	Type     int
	CardType int
	Code     string `gorm:"primary_key;type:varchar(256)"`
	Name     string
	Time     string
	Desc     string
	UpdateDt time.Time
}

type CardDataForClient struct {
	Code string
	Name string
	Time string
}

func (data *CardData) CardDataToClient() *CardDataForClient {
	tmp := new(CardDataForClient)
	tmp.Code = data.Code
	tmp.Name = data.Name
	tmp.Time = data.Time
	return tmp
}

type SecureCardData struct {
	Data map[string](map[string][]*CardData)
	Lock sync.RWMutex
}

func (data *SecureLastestCardData) String() string {
	data.Lock.Lock()
	defer data.Lock.Unlock()
	str := "\n\n==============lastest data========="
	for k, v := range data.Data {
		str += k + ":" + v.TimeStr + "\n"
	}
	return str
}
