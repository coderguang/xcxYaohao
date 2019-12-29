package config

import (
	"errors"
	"strconv"
	"strings"
	"time"
	"xcxYaohao/src/define"

	"github.com/coderguang/GameEngine_go/sglog"
	"github.com/coderguang/GameEngine_go/sgstring"
)

func GetCityName(title string) string {
	switch title {
	case define.CITY_GUANGZHOU:
		return "广州市"
	case define.CITY_SHENZHEN:
		return "深圳市"
	case define.CITY_HANGZHOU:
		return "杭州市"
	case define.CITY_TIANJIN:
		return "天津市"
	case define.CITY_HAINAN:
		return "海南市"
	case define.CITY_BEIJING:
		return "北京市"
	}
	return "未知"
}

//筛选一个页面是否应该被访问
func PageFliter(title, pageTitle, link string) bool {
	switch title {
	case define.CITY_GUANGZHOU:
		return guangzhouFliter(pageTitle)
	case define.CITY_SHENZHEN:
		return shenzhenFliter(pageTitle)
	case define.CITY_HANGZHOU:
		return hangzhouFliter(pageTitle, link)
	case define.CITY_TIANJIN:
		return tianjinFliter(pageTitle)
	case define.CITY_BEIJING:
		return beijingFliter(pageTitle)
	}
	return false
}

func guangzhouFliter(pageTitle string) bool {
	pageTxt := []string{"配置结果"}
	if !sgstring.ContainsWithOr(pageTitle, pageTxt) && !strings.Contains(pageTitle, "下一页") && !strings.Contains(pageTitle, "pdf") {
		return false
	}
	return true
}
func hangzhouFliter(pageTitle string, link string) bool {
	pageTxt := []string{"摇号结果"}
	if !sgstring.ContainsWithOr(pageTitle, pageTxt) && !strings.Contains(pageTitle, "下一页") && !strings.Contains(pageTitle, "pdf") && !strings.Contains(link, "attachment") {
		return false
	}
	return true
}

func shenzhenFliter(pageTitle string) bool {
	pageTxt := []string{"增量指标摇号结果公告", "指标配置结果"}
	if !sgstring.ContainsWithOr(pageTitle, pageTxt) && !strings.Contains(pageTitle, "下一页") && !strings.Contains(pageTitle, "pdf") {
		return false
	}
	return true
}

func tianjinFliter(pageTitle string) bool {
	pageTxt := []string{"配置结果"}
	if !sgstring.ContainsWithOr(pageTitle, pageTxt) && !strings.Contains(pageTitle, "下一页") && !strings.Contains(pageTitle, "pdf") {
		return false
	}
	return true
}

func beijingFliter(pageTitle string) bool {
	pageTxt := []string{"配置结果"}
	if !sgstring.ContainsWithOr(pageTitle, pageTxt) && !strings.Contains(pageTitle, "下一页") && !strings.Contains(pageTitle, "pdf") && !strings.Contains(pageTitle, "PDF") {
		return false
	}
	return true
}

//return data,timestr,totalnum,cardtype,membertype,error
func TxtFileFliter(title string, contents []string) (map[string]*define.CardData, string, int, int, int, error) {

	switch title {
	case define.CITY_GUANGZHOU:
		return guangzhouTxtFliter(title, contents)
	case define.CITY_SHENZHEN:
		return shenzhenTxtFliter(title, contents)
	case define.CITY_HANGZHOU:
		return hangzhouTxtFliter(title, contents)
	case define.CITY_TIANJIN:
		return tianjinTxtFliter(title, contents)
	case define.CITY_BEIJING:
		return beijingTxtFliter(title, contents)
	}
	return nil, "", 0, 0, 0, errors.New("no match title fliter")
}

func guangzhouTxtFliter(title string, contents []string) (map[string]*define.CardData, string, int, int, int, error) {

	datas := make(map[string]*define.CardData)
	timestr := ""
	totalNum := 0
	cardType := 0
	memberType := 0
	startParseData := false
	ignoreNumMath := false
	for _, v := range contents {
		if startParseData {
			if strings.Contains(v, "中签详细列表数据完成") {
				ignoreNumMath = true
			}
			strlist := strings.Split(v, " ")
			strlistex := []string{}
			for _, v := range strlist {
				if v != " " && v != "" {
					strlistex = append(strlistex, v)
				}
			}
			if len(strlistex) < 3 {
				continue
			}

			_, err := strconv.Atoi(strlistex[0])
			if err != nil {
				continue
			}

			//深圳有公司名称内有空格的，例如 美国COPPEL CORPORATION深圳代表处

			targetName := strlistex[2]
			for i := 3; i < len(strlistex); i++ {
				targetName += " " + strlistex[i]
			}

			// if len(strlistex) > 3 {
			// 	sglog.Error("not suport current format data,please check")
			// 	return datas, timestr, totalNum, cardType, memberType, errors.New("parse code name time data error,data:" + v)
			// }

			data := new(define.CardData)
			data.Title = title
			data.Time = timestr
			data.CardType = cardType
			data.Type = memberType
			codemaxlen := 50
			namemaxlen := 300
			data.Code = strlistex[1]
			data.Name = strlistex[2]
			data.UpdateDt = time.Now()

			_, err = strconv.Atoi(data.Code)
			if err != nil {
				return datas, timestr, totalNum, cardType, memberType, errors.New("code can't transform to a number")
			}

			if len(data.Code) > codemaxlen {
				sglog.Error("code len more than ", codemaxlen, ",it is ", len(data.Code), ",old code=", data.Code)
				data.Code = data.Code[0 : codemaxlen-1]
				data.Desc = "code cut "
				sglog.Error("new code=", data.Code)
			}
			if len(data.Name) > namemaxlen {
				sglog.Error("name len more than", namemaxlen, ",it is ", len(data.Name), ",old name=", data.Name)
				data.Name = data.Name[0 : namemaxlen-1]
				data.Desc += "name cut "
				sglog.Error("new name=", data.Name)
			}
			datas[data.Code] = data
		} else {
			if strings.Contains(v, "序号") {
				if "" == timestr || 0 == totalNum || 0 == memberType || 0 == cardType {
					return datas, timestr, totalNum, cardType, memberType, errors.New("head params parse error")
				}
				startParseData = true
				sglog.Info("start parse txt file detail,time:", timestr, ",num:", totalNum, "type:", memberType, ",cardType:", cardType)
				continue
			} else {
				if "" == timestr {
					if strings.Contains(v, "分期编号") {
						strlist := strings.Split(v, "：")
						if len(strlist) >= 2 {
							timestr = strlist[len(strlist)-1]
						}
					}
				}
				if 0 == totalNum {
					if strings.Contains(v, "指标总数") {
						strlist := strings.Split(v, "：")
						if len(strlist) >= 2 {
							totalNum, _ = strconv.Atoi(strlist[len(strlist)-1])
						}
					}
				}
				if 0 == memberType {
					if strings.Contains(v, "个人") {
						memberType = define.MEMBER_TYPE_PERSIONAL
					}
					if strings.Contains(v, "单位") {
						memberType = define.MEMBER_TYPE_COMPANY
					}
				}
				if 0 == cardType {
					if strings.Contains(v, "普通") {
						cardType = define.CARD_TYPE_NORMAL
					}
					newEngineTxt := []string{"新能源", "节能"}
					if sgstring.ContainsWithOr(v, newEngineTxt) {
						cardType = define.CARD_TYPE_NEW_ENGINE
					}
				}
			}
		}
	}

	if !ignoreNumMath && len(datas) != totalNum {
		return datas, timestr, totalNum, cardType, memberType, errors.New("total parse num not match,need:" + strconv.Itoa(totalNum) + ",real:" + strconv.Itoa(len(datas)))
	}
	return datas, timestr, len(datas), cardType, memberType, nil
}

func shenzhenTxtFliter(title string, contents []string) (map[string]*define.CardData, string, int, int, int, error) {

	datas := make(map[string]*define.CardData)
	timestr := ""
	totalNum := 0
	cardType := 0
	memberType := 0
	startParseData := false
	ignoreNumMath := false
	for _, v := range contents {
		if startParseData {
			if strings.Contains(v, "中签详细列表数据完成") {
				ignoreNumMath = true
			}
			strlist := strings.Split(v, " ")
			strlistex := []string{}
			for _, v := range strlist {
				if v != " " && v != "" {
					strlistex = append(strlistex, v)
				}
			}
			if len(strlistex) < 3 {
				continue
			}

			_, err := strconv.Atoi(strlistex[0])
			if err != nil {
				continue
			}

			//深圳有公司名称内有空格的，例如 美国COPPEL CORPORATION深圳代表处

			targetName := strlistex[2]
			for i := 3; i < len(strlistex); i++ {
				targetName += " " + strlistex[i]
			}

			// if len(strlistex) > 3 {
			// 	sglog.Error("not suport current format data,please check")
			// 	return datas, timestr, totalNum, cardType, memberType, errors.New("parse code name time data error,data:" + v)
			// }

			data := new(define.CardData)
			data.Title = title
			data.Time = timestr
			data.CardType = cardType
			data.Type = memberType
			codemaxlen := 50
			namemaxlen := 300
			data.Code = strlistex[1]

			_, err = strconv.Atoi(data.Code)
			if err != nil {
				return datas, timestr, totalNum, cardType, memberType, errors.New("code can't transform to a number")
			}

			//data.Name = strlistex[2]
			data.Name = targetName
			data.UpdateDt = time.Now()
			if len(data.Code) > codemaxlen {
				sglog.Error("code len more than ", codemaxlen, ",it is ", len(data.Code), ",old code=", data.Code)
				data.Code = data.Code[0 : codemaxlen-1]
				data.Desc = "code cut "
				sglog.Error("new code=", data.Code)
			}
			if len(data.Name) > namemaxlen {
				sglog.Error("name len more than", namemaxlen, ",it is ", len(data.Name), ",old name=", data.Name)
				data.Name = data.Name[0 : namemaxlen-1]
				data.Desc += "name cut "
				sglog.Error("new name=%s", data.Name)
			}
			datas[data.Code] = data
		} else {
			if strings.Contains(v, "序号") {
				if "" == timestr || 0 == totalNum || 0 == memberType || 0 == cardType {
					return datas, timestr, totalNum, cardType, memberType, errors.New("head params parse error")
				}
				startParseData = true
				sglog.Info("start parse txt file detail,time:", timestr, ",num:", totalNum, "type:", memberType, ",cardType:", cardType)
				continue
			} else {
				if "" == timestr {
					if strings.Contains(v, "本期编号") {
						strlist := strings.Split(v, "：")
						if len(strlist) >= 2 {
							timestr = strlist[len(strlist)-1]
						}
					}
				}
				if 0 == totalNum {
					if strings.Contains(v, "指标总数") {
						strlist := strings.Split(v, "：")
						if len(strlist) >= 2 {
							totalNum, _ = strconv.Atoi(strlist[len(strlist)-1])
						}
					}
				}
				if 0 == memberType {
					if strings.Contains(v, "个人") {
						memberType = define.MEMBER_TYPE_PERSIONAL
					}
					if strings.Contains(v, "单位") {
						memberType = define.MEMBER_TYPE_COMPANY
					}
				}
				if 0 == cardType {
					if strings.Contains(v, "普通") {
						cardType = define.CARD_TYPE_NORMAL
					}
					newEngineTxt := []string{"新能源", "电动"}
					if sgstring.ContainsWithOr(v, newEngineTxt) {
						cardType = define.CARD_TYPE_NEW_ENGINE
					}
				}
			}
		}
	}

	if !ignoreNumMath && len(datas) != totalNum {
		return datas, timestr, totalNum, cardType, memberType, errors.New("total parse num not match,need:" + strconv.Itoa(totalNum) + ",real:" + strconv.Itoa(len(datas)))
	}
	return datas, timestr, len(datas), cardType, memberType, nil
}

func hangzhouTxtFliter(title string, contents []string) (map[string]*define.CardData, string, int, int, int, error) {

	datas := make(map[string]*define.CardData)
	timestr := ""
	totalNum := 0
	cardType := 0
	memberType := 0
	startParseData := false
	ignoreNumMath := false
	for _, v := range contents {
		if startParseData {
			if strings.Contains(v, "中签详细列表数据完成") {
				ignoreNumMath = true
			}
			strlist := strings.Split(v, " ")
			strlistex := []string{}
			for _, v := range strlist {
				if v != " " && v != "" {
					strlistex = append(strlistex, v)
				}
			}
			if len(strlistex) < 3 {
				continue
			}

			_, err := strconv.Atoi(strlistex[0])
			if err != nil {
				continue
			}

			//深圳有公司名称内有空格的，例如 美国COPPEL CORPORATION深圳代表处

			targetName := strlistex[2]
			for i := 3; i < len(strlistex); i++ {
				targetName += " " + strlistex[i]
			}

			// if len(strlistex) > 3 {
			// 	sglog.Error("not suport current format data,please check")
			// 	return datas, timestr, totalNum, cardType, memberType, errors.New("parse code name time data error,data:" + v)
			// }

			data := new(define.CardData)
			data.Title = title
			data.Time = timestr
			data.CardType = cardType
			data.Type = memberType
			codemaxlen := 50
			namemaxlen := 300
			data.Code = strlistex[1]
			data.Name = strlistex[2]
			data.UpdateDt = time.Now()

			_, err = strconv.Atoi(data.Code)
			if err != nil {
				return datas, timestr, totalNum, cardType, memberType, errors.New("code can't transform to a number")
			}

			if len(data.Code) > codemaxlen {
				sglog.Error("code len more than ", codemaxlen, ",it is ", len(data.Code), ",old code=", data.Code)
				data.Code = data.Code[0 : codemaxlen-1]
				data.Desc = "code cut "
				sglog.Error("new code=", data.Code)
			}
			if len(data.Name) > namemaxlen {
				sglog.Error("name len more than", namemaxlen, ",it is ", len(data.Name), ",old name=", data.Name)
				data.Name = data.Name[0 : namemaxlen-1]
				data.Desc += "name cut "
				sglog.Error("new name=%s", data.Name)
			}
			datas[data.Code] = data
		} else {
			if strings.Contains(v, "序号") {
				if "" == timestr || 0 == totalNum || 0 == memberType || 0 == cardType {
					return datas, timestr, totalNum, cardType, memberType, errors.New("head params parse error")
				}
				startParseData = true
				sglog.Info("start parse txt file detail,time:", timestr, ",num:", totalNum, "type:", memberType, ",cardType:", cardType)
				continue
			} else {
				if "" == timestr {
					if strings.Contains(v, "分期编号") {
						strlist := strings.Split(v, "：")
						if len(strlist) >= 2 {
							timestr = strlist[len(strlist)-1]
						}
					}
				}
				if 0 == totalNum {
					if strings.Contains(v, "指标总数") {
						strlist := strings.Split(v, "：")
						if len(strlist) >= 2 {
							totalNum, _ = strconv.Atoi(strlist[len(strlist)-1])
						}
					}
				}
				if 0 == memberType {
					if strings.Contains(v, "个人") {
						memberType = define.MEMBER_TYPE_PERSIONAL
					}
					if strings.Contains(v, "单位") {
						memberType = define.MEMBER_TYPE_COMPANY
					}
				}
				if 0 == cardType {
					cardType = define.CARD_TYPE_NORMAL
				}
			}
		}
	}

	if !ignoreNumMath && len(datas) != totalNum {
		return datas, timestr, totalNum, cardType, memberType, errors.New("total parse num not match,need:" + strconv.Itoa(totalNum) + ",real:" + strconv.Itoa(len(datas)))
	}
	return datas, timestr, len(datas), cardType, memberType, nil
}

func tianjinTxtFliter(title string, contents []string) (map[string]*define.CardData, string, int, int, int, error) {

	datas := make(map[string]*define.CardData)
	timestr := ""
	totalNum := 0
	cardType := 0
	memberType := 0
	startParseData := false
	ignoreNumMath := false
	for _, v := range contents {
		if startParseData {
			if strings.Contains(v, "中签详细列表数据完成") {
				ignoreNumMath = true
			}
			strlist := strings.Split(v, " ")
			strlistex := []string{}
			for _, v := range strlist {
				if v != " " && v != "" {
					strlistex = append(strlistex, v)
				}
			}
			if len(strlistex) < 3 {
				continue
			}

			_, err := strconv.Atoi(strlistex[0])
			if err != nil {
				continue
			}

			//深圳有公司名称内有空格的，例如 美国COPPEL CORPORATION深圳代表处

			targetName := strlistex[2]
			for i := 3; i < len(strlistex); i++ {
				targetName += " " + strlistex[i]
			}

			// if len(strlistex) > 3 {
			// 	sglog.Error("not suport current format data,please check")
			// 	return datas, timestr, totalNum, cardType, memberType, errors.New("parse code name time data error,data:" + v)
			// }

			data := new(define.CardData)
			data.Title = title
			data.Time = timestr
			data.CardType = cardType
			data.Type = memberType
			codemaxlen := 50
			namemaxlen := 300
			data.Code = strlistex[1]
			data.Name = strlistex[2]

			_, err = strconv.Atoi(data.Code)
			if err != nil {
				return datas, timestr, totalNum, cardType, memberType, errors.New("code can't transform to a number")
			}

			data.UpdateDt = time.Now()
			if len(data.Code) > codemaxlen {
				sglog.Error("code len more than ", codemaxlen, ",it is ", len(data.Code), ",old code=", data.Code)
				data.Code = data.Code[0 : codemaxlen-1]
				data.Desc = "code cut "
				sglog.Error("new code=", data.Code)
			}
			if len(data.Name) > namemaxlen {
				sglog.Error("name len more than", namemaxlen, ",it is ", len(data.Name), ",old name=", data.Name)
				data.Name = data.Name[0 : namemaxlen-1]
				data.Desc += "name cut "
				sglog.Error("new name=%s", data.Name)
			}
			datas[data.Code] = data
		} else {
			if strings.Contains(v, "序号") {
				if "" == timestr || 0 == totalNum || 0 == memberType || 0 == cardType {
					return datas, timestr, totalNum, cardType, memberType, errors.New("head params parse error")
				}
				startParseData = true
				sglog.Info("start parse txt file detail,time:", timestr, ",num:", totalNum, "type:", memberType, ",cardType:", cardType)
				continue
			} else {
				if "" == timestr {
					if strings.Contains(v, "分期编号") {
						strlist := strings.Split(v, "：")
						if len(strlist) >= 2 {
							timestr = strlist[len(strlist)-1]
						}
					}
				}
				if 0 == totalNum {
					if strings.Contains(v, "指标总数") {
						strlist := strings.Split(v, "：")
						if len(strlist) >= 2 {
							totalNum, _ = strconv.Atoi(strlist[len(strlist)-1])
						}
					}
				}
				if 0 == memberType {
					if strings.Contains(v, "个人") {
						memberType = define.MEMBER_TYPE_PERSIONAL
					}
					if strings.Contains(v, "单位") {
						memberType = define.MEMBER_TYPE_COMPANY
					}
				}
				if 0 == cardType {
					if strings.Contains(v, "普通") {
						cardType = define.CARD_TYPE_NORMAL
					}
					newEngineTxt := []string{"新能源", "节能"}
					if sgstring.ContainsWithOr(v, newEngineTxt) {
						cardType = define.CARD_TYPE_NEW_ENGINE
					}
				}
			}
		}
	}

	if !ignoreNumMath && len(datas) != totalNum {
		return datas, timestr, totalNum, cardType, memberType, errors.New("total parse num not match,need:" + strconv.Itoa(totalNum) + ",real:" + strconv.Itoa(len(datas)))
	}
	return datas, timestr, len(datas), cardType, memberType, nil
}

func beijingTxtFliter(title string, contents []string) (map[string]*define.CardData, string, int, int, int, error) {

	datas := make(map[string]*define.CardData)
	timestr := ""
	totalNum := 0
	cardType := 0
	memberType := 0
	startParseData := false
	ignoreNumMath := false
	isNormalOrder := false
	isTurnTimes := false
	isOrderNum := false
	isHadName := false
	for _, v := range contents {
		if startParseData {
			if strings.Contains(v, "中签详细列表数据完成") {
				ignoreNumMath = true
			}
			strlist := strings.Split(v, " ")
			strlistex := []string{}
			for _, v := range strlist {
				if v != " " && v != "" {
					strlistex = append(strlistex, v)
				}
			}
			if len(strlistex) < 3 {
				continue
			}

			_, err := strconv.Atoi(strlistex[0])
			if err != nil {
				continue
			}

			data := new(define.CardData)
			data.Title = title
			data.Time = timestr
			data.CardType = cardType
			data.Type = memberType
			codemaxlen := 50
			namemaxlen := 300
			if isTurnTimes { //轮候序号 轮候时间 申请编码
				if len(strlistex) != 4 {
					return datas, timestr, totalNum, cardType, memberType, errors.New("turn format error,len not equal 4")
				}
				data.Code = strlistex[3]
				data.TurnsTime = strlistex[1]
				// orderN, err := strconv.Atoi(strlistex[0])
				// if err != nil {
				// 	return datas, timestr, totalNum, cardType, memberType, errors.New("turn format error,order can't trans to num")
				// }
				// data.OrderNum = orderN
			} else {
				if isOrderNum && isHadName { //序号 摇号基数序号 申请编码 名称
					if len(strlistex) != 4 {
						return datas, timestr, totalNum, cardType, memberType, errors.New("order code name error,len not equal 4")
					}
					data.Code = strlistex[2]
					data.Name = strlistex[3]
					orderN, err := strconv.Atoi(strlistex[1])
					if err != nil {
						return datas, timestr, totalNum, cardType, memberType, errors.New("order code name  error,order can't trans to num")
					}
					data.OrderNum = orderN
				} else if isOrderNum { //序号 摇号基数序号 申请编码
					if len(strlistex) != 3 {
						return datas, timestr, totalNum, cardType, memberType, errors.New("order code error,len not equal 4")
					}
					data.Code = strlistex[2]
					orderN, err := strconv.Atoi(strlistex[1])
					if err != nil {
						return datas, timestr, totalNum, cardType, memberType, errors.New("order code error,order can't trans to num")
					}
					data.OrderNum = orderN
				} else if isNormalOrder { //序号 申请编码 名称
					if len(strlistex) < 3 {
						return datas, timestr, totalNum, cardType, memberType, errors.New("order code error,len not equal 4")
					}
					targetName := ""
					for i := 2; i < len(strlistex); i++ {
						targetName += strlistex[i]
					}

					data.Code = strlistex[1]
					data.Name = targetName
				} else {
					return datas, timestr, totalNum, cardType, memberType, errors.New("bejing unknow format error,please check")
				}
			}
			data.UpdateDt = time.Now()

			_, err = strconv.Atoi(data.Code)
			if err != nil {
				return datas, timestr, totalNum, cardType, memberType, errors.New("code can't transform to a number")
			}

			if len(data.Code) != 13 {
				return datas, timestr, totalNum, cardType, memberType, errors.New("code len not 13")
			}

			if len(data.Code) > codemaxlen {
				sglog.Error("code len more than ", codemaxlen, ",it is ", len(data.Code), ",old code=", data.Code)
				data.Code = data.Code[0 : codemaxlen-1]
				data.Desc = "code cut "
				sglog.Error("new code=", data.Code)
			}
			if len(data.Name) > namemaxlen {
				sglog.Error("name len more than", namemaxlen, ",it is ", len(data.Name), ",old name=", data.Name)
				data.Name = data.Name[0 : namemaxlen-1]
				data.Desc += "name cut "
				sglog.Error("new name=", data.Name)
			}
			datas[data.Code] = data
		} else {
			if strings.Contains(v, "序号") {
				if "" == timestr || 0 == totalNum || 0 == memberType {
					return datas, timestr, totalNum, cardType, memberType, errors.New("head params parse error")
				}
				if 0 == cardType {
					cardType = define.CARD_TYPE_NORMAL
					sglog.Info("title:", title, ",time:", timestr, ",no cardType")
				}
				if sgstring.ContainsWithAnd(v, []string{"申请编码"}) {
					if sgstring.ContainsWithOr(v, []string{"姓名", "名称"}) {
						isNormalOrder = true
					}
				}
				if strings.Contains(v, "轮候时间") {
					isTurnTimes = true
				}
				if strings.Contains(v, "摇号基数") {
					isOrderNum = true
				}
				if strings.Contains(v, "名称") {
					isHadName = true
				}

				startParseData = true
				sglog.Info("start parse txt file detail,time:", timestr, ",num:", totalNum, "type:", memberType, ",cardType:", cardType)
				continue
			} else {
				if "" == timestr {
					timelist := []string{"分期编号", "期号"}
					if sgstring.ContainsWithOr(v, timelist) {
						strlist := strings.Split(v, "：")
						if len(strlist) >= 2 {
							timestr = strlist[len(strlist)-1]
						}
					}
				}
				if 0 == totalNum {
					if sgstring.ContainsWithOr(v, []string{"指标总数", "指标配置总数"}) {
						strlist := strings.Split(v, "：")
						if len(strlist) >= 2 {
							totalNum, _ = strconv.Atoi(strlist[len(strlist)-1])
						}
					}
				}
				if 0 == memberType {
					if strings.Contains(v, "个人") {
						memberType = define.MEMBER_TYPE_PERSIONAL
					}
					if strings.Contains(v, "单位") {
						memberType = define.MEMBER_TYPE_COMPANY
					}
				}
				if 0 == cardType {
					if strings.Contains(v, "普通") {
						cardType = define.CARD_TYPE_NORMAL
					}
					newEngineTxt := []string{"新能源", "节能"}
					if sgstring.ContainsWithOr(v, newEngineTxt) {
						cardType = define.CARD_TYPE_NEW_ENGINE
					}
				}
			}
		}
	}

	if !ignoreNumMath && len(datas) != totalNum {
		return datas, timestr, totalNum, cardType, memberType, errors.New("total parse num not match,need:" + strconv.Itoa(totalNum) + ",real:" + strconv.Itoa(len(datas)))
	}
	return datas, timestr, len(datas), cardType, memberType, nil
}
