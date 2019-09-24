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

//筛选一个页面是否应该被访问
func PageFliter(title, pageTitle string) bool {
	switch title {
	case "guangzhou":
		return guangzhouFliter(pageTitle)
	case "shenzhen":
		return shenzhenFliter(pageTitle)
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

func shenzhenFliter(pageTitle string) bool {
	pageTxt := []string{"增量指标摇号结果公告", "指标配置结果"}
	if !sgstring.ContainsWithOr(pageTitle, pageTxt) && !strings.Contains(pageTitle, "下一页") && !strings.Contains(pageTitle, "pdf") {
		return false
	}
	return true
}

//return data,timestr,totalnum,cardtype,membertype,error
func TxtFileFliter(title string, contents []string) (map[string]*define.CardData, string, int, int, int, error) {

	switch title {
	case "guangzhou":
		return guangzhouTxtFliter(title, contents)
	case "shenzhen":
		return shenzhenTxtFliter(title, contents)
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
			if len(strlistex) > 3 {
				sglog.Error("not suport current format data,please check")
				return datas, timestr, totalNum, cardType, memberType, errors.New("parse code name time data error,data:" + v)
			}

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
						cardType = define.CARD_TYPE_PERSION
					}
					newEngineTxt := []string{"新能源", "节能"}
					if sgstring.ContainsWithOr(v, newEngineTxt) {
						cardType = define.CARD_TYPE_COMPANY
					}
				}
			}
		}
	}

	if !ignoreNumMath && len(datas) != totalNum {
		return datas, timestr, totalNum, cardType, memberType, errors.New("total parse num not match")
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
			if len(strlistex) > 3 {
				sglog.Error("not suport current format data,please check")
				return datas, timestr, totalNum, cardType, memberType, errors.New("parse code name time data error,data:" + v)
			}

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
						cardType = define.CARD_TYPE_PERSION
					}
					newEngineTxt := []string{"新能源", "电动"}
					if sgstring.ContainsWithOr(v, newEngineTxt) {
						cardType = define.CARD_TYPE_COMPANY
					}
				}
			}
		}
	}

	if !ignoreNumMath && len(datas) != totalNum {
		return datas, timestr, totalNum, cardType, memberType, errors.New("total parse num not match")
	}
	return datas, timestr, len(datas), cardType, memberType, nil
}
