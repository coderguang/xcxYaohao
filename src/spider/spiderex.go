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

	"github.com/PuerkitoBio/goquery"
	"github.com/coderguang/GameEngine_go/sglog"
	"github.com/coderguang/GameEngine_go/sgstring"
	"github.com/coderguang/GameEngine_go/sgthread"
	"github.com/coderguang/GameEngine_go/sgtime"
)

func StartSpiderEx(title string, startDt time.Time, index string, isPersonal bool) {

	searchPage := 1
	for {

		searchTime := sgtime.YMString(sgtime.TransfromTimeToDateTime(startDt))

		sglog.Info(title, "start spider date:", searchTime, ",isPersonal:", isPersonal)

		dataMap := make(map[string]string)

		searchPageStr := strconv.Itoa(searchPage)
		ignoreMap := make(map[string]string)
		totalPage, totalNum, err := dataSpider(title, index, searchTime, searchPageStr, ignoreMap)
		if err != nil {
			sglog.Error("search err,", searchTime, searchPageStr, err)
			continue
		}
		deleteRecord := 0
		totalPage = 1
		for i := 1; i <= totalPage; i++ {
			sgthread.SleepByMillSecond(50)
			tmpDataMap := make(map[string]string)
			_, _, err = dataSpider(title, index, searchTime, strconv.Itoa(i), tmpDataMap)
			if err != nil {
				sglog.Error(title, "find spider err,", searchTime, i, err)
				i--
				continue
			}
			if len(tmpDataMap) != 16 {
				sglog.Debug("page size not 16,", searchTime, i, ",size is ", len(tmpDataMap))
			}
			sglog.Debug("page:", i, "size:", len(tmpDataMap), ",totalsize:", len(dataMap))
			for k, v := range tmpDataMap {
				if vv, ok := dataMap[k]; ok {
					sglog.Info(title, "duplicate data,", k, v, vv)
					if v == vv {
						deleteRecord++
					}
				}
				dataMap[k] = v
			}
		}
		if totalNum != len(dataMap)+deleteRecord {
			sglog.Error(title, "end search ", searchTime, ",needSize:", totalNum, ",real size:", len(dataMap))
			sgthread.DelayExit(2)
			continue
		} else {
			sglog.Info(title, "end search ok", searchTime, ",needSize:", totalNum, ",real size:", len(dataMap))
			now := time.Now()
			cardDataMap := make(map[string]*define.CardData)
			for k, v := range dataMap {
				tmp := new(define.CardData)
				tmp.Title = title
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
				tmp.Desc = "by " + title + " special"
				tmp.UpdateDt = now

				if data.IsDataExist(title, k) {
					sglog.Error("data already exist,time:", searchTime, "code:", k)
				} else {
					cardDataMap[k] = tmp
				}
			}
			data.AddCardData(cardDataMap)

			go func(updates map[string]*define.CardData) {
				nowEx := sgtime.New()
				updateDbNum := 0
				for _, v := range updates {
					if err = db.UpdateCardData(v); err != nil {
						sglog.Error(title, "update data to db error,title", v.Title, "code", v.Code, err)
					} else {
						updateDbNum++
					}
				}
				endEx := sgtime.New()
				sglog.Info(title, " update card data in databases size:", updateDbNum, ",use time:", (sgtime.GetTotalSecond(endEx) - sgtime.GetTotalSecond(nowEx)))

			}(cardDataMap)

		}
		if searchTime == "201910" {
			break
		}
		startDt = startDt.AddDate(0, 1, 0)
	}

	sglog.Info("spider ", title, " old data ok,isPersion:", isPersonal)

}

func dataSpider(title string, index string, timestr string, page string, dataMap map[string]string) (int, int, error) {
	totalPage := 0
	totalNum := 0
	params := "pageNo=" + page + "&issueNumber=" + timestr + "&applyCode= "
	resp, err := http.Post(index,
		"application/x-www-form-urlencoded",
		strings.NewReader(params))

	if err != nil {
		sglog.Error(title, " special err,", err)
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
					if httpHandle.CheckCodeValid(title, code) {
						if v, ok := dataMap[code]; ok {
							sglog.Error(title, "duplicate code:", code, ",name:", name, ",time:", timestr, ",oldname:", v, ",page:", page)
						} else {
							dataMap[code] = name
						}
					} else {
						sglog.Error(title, "code not valid,code:", code)
					}
				}
			})
		} else if tmpS.Size() == 3 {
			tr.Find("td").Each(func(ix int, td *goquery.Selection) {
				sglog.Info(ix, ":", td.Text())
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
