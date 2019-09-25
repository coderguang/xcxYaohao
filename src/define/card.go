package define

import (
	"sync"
	"time"
)

const (
	CARD_TYPE_PERSION     = 1
	CARD_TYPE_COMPANY     = 2
	MEMBER_TYPE_PERSIONAL = 1
	MEMBER_TYPE_COMPANY   = 2
)

type SLastestCardData struct {
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
		if CARD_TYPE_PERSION == cardType {
			data.PersonalNormalUpdate = true
		} else {
			data.PersonalJieNengUpdate = true
		}
	} else {
		if CARD_TYPE_PERSION == cardType {
			data.CompanyNormalUpdate = true
		} else {
			data.CompanyJieNengUpdate = true
		}
	}
}

func (data *SLastestCardData) IsAllCardInfoUpdate() bool {
	if data.PersonalJieNengUpdate && data.PersonalNormalUpdate && data.CompanyJieNengUpdate && data.CompanyNormalUpdate {
		return true
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
