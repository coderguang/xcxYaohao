package data

import (
	"sort"
	"xcxYaohao/src/define"

	"github.com/coderguang/GameEngine_go/sglog"
	"github.com/mohae/deepcopy"
)

var (
	globalCardData       *define.SecureCardData
	globalCardDataByName *define.SecureCardData
	globalLastestData    *define.SecureLastestCardData
)

func init() {
	globalCardData = new(define.SecureCardData)
	globalCardData.Data = make(map[string](map[string][]*define.CardData))

	globalCardDataByName = new(define.SecureCardData)
	globalCardDataByName.Data = make(map[string](map[string][]*define.CardData))

	globalLastestData = new(define.SecureLastestCardData)
	globalLastestData.Data = make(map[string]*define.SLastestCardData)
}

func GetLastestCardInfo(title string) define.SLastestCardData {
	return define.SLastestCardData{}
}

func IsAllCardUpdate(title string, dt string) bool {
	return true
}

func InitCardDataFromDb(datas []define.CardData) {

	globalCardData.Lock.Lock()
	defer globalCardData.Lock.Unlock()
	for _, v := range datas {
		UpdateLastestInfo(v.Title, v.CardType, v.Type, v.Time)
		tmp := deepcopy.Copy(v)
		tmpV, ok := tmp.(define.CardData)
		if !ok {
			sglog.Error("deepcopy error,title:", v.Title, ",code:", v.Code)
			continue
		}

		if cityMap, ok := globalCardData.Data[v.Title]; ok {
			if codeList, ok := cityMap[v.Code]; ok {
				codeList = append(codeList, &tmpV)
			} else {
				tmp := []*define.CardData{&tmpV}
				cityMap[v.Code] = tmp
			}
		} else {
			globalCardData.Data[v.Title] = make(map[string][]*define.CardData)
			tmp := []*define.CardData{&tmpV}
			globalCardData.Data[v.Title][v.Code] = tmp
		}
	}

	for _, v := range globalCardData.Data {
		for _, vv := range v {
			sort.Slice(vv, func(i, j int) bool {
				return vv[i].Time > vv[j].Time
			})
		}
	}

	globalCardDataByName.Lock.Lock()
	defer globalCardDataByName.Lock.Unlock()
	for _, v := range datas {
		tmp := deepcopy.Copy(v)
		tmpV, ok := tmp.(define.CardData)
		if !ok {
			sglog.Error("deepcopy by name error,title:", v.Title, ",code:", v.Code)
			continue
		}
		if cityMap, ok := globalCardDataByName.Data[v.Title]; ok {
			if namelist, ok := cityMap[v.Name]; ok {
				namelist = append(namelist, &tmpV)
			} else {
				tmp := []*define.CardData{&tmpV}
				cityMap[v.Code] = tmp
			}
		} else {
			globalCardDataByName.Data[v.Title] = make(map[string][]*define.CardData)
			tmp := []*define.CardData{&tmpV}
			globalCardDataByName.Data[v.Title][v.Name] = tmp
		}
	}

	for _, v := range globalCardDataByName.Data {
		for _, vv := range v {
			sort.Slice(vv, func(i, j int) bool {
				return vv[i].Time > vv[j].Time
			})
		}
	}

}

func UpdateLastestInfo(title string, cardType int, memberType int, timestr string) {
	globalLastestData.Lock.Lock()
	defer globalLastestData.Lock.Unlock()

	if v, ok := globalLastestData.Data[title]; ok {
		if v.TimeStr < timestr {
			v.Reset()
			v.TimeStr = timestr
		} else if v.TimeStr == timestr {
			v.UpdateLastestInfo(cardType, memberType)
		}
		return
	}
	tmp := new(define.SLastestCardData)
	tmp.Reset()
	tmp.TimeStr = timestr
	tmp.UpdateLastestInfo(cardType, memberType)
	globalLastestData.Data[title] = tmp
}

func IsDataExist(title string, code string) bool {
	globalCardData.Lock.Lock()
	defer globalCardData.Lock.Unlock()
	if cityMap, ok := globalCardData.Data[title]; ok {
		if _, ok = cityMap[code]; ok {
			return true
		}
	}
	return false
}

func AddCardData(datas map[string]*define.CardData) {
	globalCardData.Lock.Lock()
	defer globalCardData.Lock.Unlock()

	for _, v := range datas {
		if cityMap, ok := globalCardData.Data[v.Title]; ok {
			if namelist, ok := cityMap[v.Code]; ok {
				namelist = append(namelist, v)
			} else {
				tmp := []*define.CardData{v}
				cityMap[v.Code] = tmp
			}
		} else {
			globalCardData.Data[v.Title] = make(map[string][]*define.CardData)
			tmp := []*define.CardData{v}
			globalCardData.Data[v.Title][v.Code] = tmp
		}
	}
	globalCardDataByName.Lock.Lock()
	defer globalCardDataByName.Lock.Unlock()

	for _, v := range datas {
		if cityMap, ok := globalCardDataByName.Data[v.Title]; ok {
			if namelist, ok := cityMap[v.Name]; ok {
				namelist = append(namelist, v)
			} else {
				tmp := []*define.CardData{v}
				cityMap[v.Code] = tmp
			}
		} else {
			globalCardDataByName.Data[v.Title] = make(map[string][]*define.CardData)
			tmp := []*define.CardData{v}
			globalCardDataByName.Data[v.Title][v.Name] = tmp
		}
	}

}

func ShowLastestInfo(cmd []string) {
	sglog.Debug(globalLastestData)
}
