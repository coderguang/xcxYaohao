package spider

import (
	"net/http"
	"strconv"
	"strings"
	"time"
	"xcxYaohao/src/data"
	"xcxYaohao/src/db"
	"xcxYaohao/src/define"
	"xcxYaohao/src/httpHandle"

	"github.com/coderguang/GameEngine_go/sgthread"

	"github.com/coderguang/GameEngine_go/sgtime"

	"github.com/coderguang/GameEngine_go/sgstring"

	"github.com/PuerkitoBio/goquery"
	"github.com/coderguang/GameEngine_go/sglog"
)

func TianjinOldDataSpider(cmd []string) {
	unitIndex := "http://apply.xkctk.jtys.tj.gov.cn/apply/norm/unitQuery.html"
	persionIndex := "http://apply.xkctk.jtys.tj.gov.cn/apply/norm/personQuery.html"
	StartSpiderTianJinOldData(unitIndex, false)
	StartSpiderTianJinOldData(persionIndex, true)

}

func StartSpiderTianJinOldData(index string, isPersonal bool) {

	startDt, err := sgtime.ParseInLocation(sgtime.FORMAT_TIME_NORMAL, "2014-02-01 00:00:00")
	if err != nil {
		sglog.Error("parse startDt error,", err)
	}

	searchPage := 1
	for {

		searchTime := sgtime.YMString(sgtime.TransfromTimeToDateTime(startDt))

		sglog.Info("start spider date:", searchTime, ",isPersonal:", isPersonal)

		dataMap := make(map[string]string)

		searchPageStr := strconv.Itoa(searchPage)
		totalPage, totalNum, err := dataSpider(index, searchTime, searchPageStr, dataMap)
		if err != nil {
			sglog.Error("search err,", searchTime, searchPageStr, err)
			continue
		}

		for i := 1; i <= totalPage; i++ {
			sgthread.SleepByMillSecond(200)
			_, _, err = dataSpider(index, searchTime, strconv.Itoa(i), dataMap)
			if err != nil {
				i--
				continue
			}
		}
		if totalNum != len(dataMap) {
			sglog.Error("end search ", searchTime, ",needSize:", totalNum, ",real size:", len(dataMap))
			continue
		} else {
			sglog.Info("end search ok", searchTime, ",needSize:", totalNum, ",real size:", len(dataMap))
			now := time.Now()
			cardDataMap := make(map[string]*define.CardData)
			for k, v := range dataMap {
				tmp := new(define.CardData)
				tmp.Title = define.CITY_TIANJIN
				tmp.Type = define.CARD_TYPE_NORMAL
				tmp.CardType = define.CARD_TYPE_NORMAL
				if isPersonal {
					tmp.Type = define.MEMBER_TYPE_PERSIONAL
				} else {
					tmp.Type = define.MEMBER_TYPE_COMPANY
				}
				tmp.Code = k
				tmp.Name = v
				tmp.Time = searchTime
				tmp.Desc = "by tianjin special"
				tmp.UpdateDt = now

				if data.IsDataExist(define.CITY_TIANJIN, k) {
					sglog.Error("data already exist,time:", searchTime, "code:", k)
				} else {
					err = db.UpdateCardData(tmp)
					if err != nil {
						sglog.Error("tianjin speical update to db error,time:", searchTime, isPersonal, err)
					}
					cardDataMap[k] = tmp
				}
			}
			data.AddCardData(cardDataMap)

		}
		if searchTime == "201910" {
			break
		}
		startDt = startDt.AddDate(0, 1, 0)
	}

	sglog.Info("spider tianjin old data ok,isPersion:", isPersonal)

}

func dataSpider(index string, timestr string, page string, dataMap map[string]string) (int, int, error) {
	totalPage := 0
	totalNum := 0
	params := "pageNo=" + page + "&issueNumber=" + timestr + "&applyCode= "
	resp, err := http.Post(index,
		"application/x-www-form-urlencoded",
		strings.NewReader(params))

	if err != nil {
		sglog.Error("tianjin special err,", err)
		return totalPage, totalNum, err
	}

	doc, err := goquery.NewDocumentFromResponse(resp)

	if err != nil {
		sglog.Error("docs error")
		return totalPage, totalNum, err
	}
	doc.Find("table tr").Each(func(_ int, tr *goquery.Selection) {
		// for each <tr> found, find the <td>s inside
		// ix is the index
		tmpS := tr.Find("td")

		if tmpS.Size() == 2 {
			code := ""
			tr.Find("td").Each(func(ix int, td *goquery.Selection) {
				if ix == 0 {
					code = td.Text()
				} else {
					name := td.Text()
					//sglog.Info(code, ":", name)
					if httpHandle.CheckCodeValid(define.CITY_TIANJIN, code) {
						dataMap[code] = name
					}
				}
			})
		} else if tmpS.Size() == 3 {
			tr.Find("td").Each(func(ix int, td *goquery.Selection) {
				//sglog.Info(ix, ":", td.Text())
				if 2 == ix {
					str := td.Text()
					if sgstring.ContainsWithAnd(str, []string{"共", "/", "页", "条"}) {
						str = strings.Replace(str, "共", "", -1)
						str = strings.Replace(str, "页", "", -1)
						str = strings.Replace(str, "条", "", -1)
						strlist := strings.Split(str, "/")
						if 2 == len(strlist) {
							totalPage, err = strconv.Atoi(strlist[0])
							if err != nil {
								sglog.Error("tranform error,", td.Text(), err)
							}
							totalNum, err = strconv.Atoi(strlist[1])
							if err != nil {
								sglog.Error("tranform error,", td.Text(), err)
							}
							//sglog.Info(td.Text(), "parse: page", totalPage, ",num:", totalNum)
						}
					}
				}
			})

		}
	})
	return totalPage, totalNum, nil
}
