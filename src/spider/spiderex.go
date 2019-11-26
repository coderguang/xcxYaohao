package spider

import (
	"net/http"
	"strconv"
	"strings"
	"time"
	"xcxYaohao/src/cache"
	"xcxYaohao/src/data"
	"xcxYaohao/src/db"
	"xcxYaohao/src/define"
	"xcxYaohao/src/httpHandle"
	"xcxYaohao/src/notice"

	"github.com/PuerkitoBio/goquery"
	"github.com/coderguang/GameEngine_go/sglog"
	"github.com/coderguang/GameEngine_go/sgstring"
	"github.com/coderguang/GameEngine_go/sgthread"
	"github.com/coderguang/GameEngine_go/sgtime"
)

func StartSpiderEx(title string, startDt time.Time, index string, isPersonal bool) {

	searchPage := 1
	sleepTime := define.SPIDER_SLEEP_TIME
	timeInt := time.Duration(define.SPIDER_TIME_INT) * time.Second
	for {

		searchTime := sgtime.YMString(sgtime.TransfromTimeToDateTime(startDt))
		urlTips := searchTime
		if isPersonal {
			urlTips += "personal"
		} else {
			urlTips += "company"
		}
		if !data.NeedDownloadFile(title, urlTips) {
			startDt = startDt.AddDate(0, 1, 0)
			sglog.Info(title, "had already download ", searchTime, ",isPersonal:", isPersonal)
			continue
		}

		nowTime := time.Now()
		curTimeStr := sgtime.YearString(&nowTime) + sgtime.MonthString(&nowTime)
		if searchTime == curTimeStr {
			normalTime := time.Date(nowTime.Year(), nowTime.Month(), define.SPIDER_START_DAY, define.SPIDER_START_HOUR, 0, 0, 0, nowTime.Location())
			if nowTime.Before(normalTime) {
				timeInt = normalTime.Sub(nowTime)
			} else {
				hour := time.Now().Hour()
				if hour < define.SPIDER_START_HOUR {
					nextRun := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), define.SPIDER_START_HOUR, 0, 0, 0, nowTime.Location())
					timeInt = nextRun.Sub(nowTime)
				} else if hour > define.SPIDER_END_HOUR {
					nextRun := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 23, 59, 59, 0, nowTime.Location())
					timeInt = nextRun.Sub(nowTime)
				}
			}
			sleepTime = int(timeInt/time.Second) + 1
			sglog.Info(title, "ex data collection now in sleep,will run after ", sleepTime, "s,", nowTime.Add(timeInt), "isPerson:", isPersonal)

			sgthread.SleepBySecond(sleepTime)
		} else if searchTime > curTimeStr {
			nextDt := nowTime.AddDate(0, 1, 0)
			nextSearchDt := time.Date(nextDt.Year(), nextDt.Month(), define.SPIDER_START_DAY, define.SPIDER_START_HOUR, 0, 0, 0, nextDt.Location())
			timeInt := nextSearchDt.Sub(nowTime)
			sleepTime = int(timeInt/time.Second) + 1
			sglog.Info(title, "ex data collection now in sleep,will run after ", sleepTime, "s,", nextSearchDt, "isPerson:", isPersonal)
			sgthread.SleepBySecond(sleepTime)
		}

		sglog.Info(title, "start spider date:", searchTime, ",isPersonal:", isPersonal)

		if title == define.CITY_HAINAN {
			if !HainanIsShowResultInWebIndex(searchTime) {
				sglog.Debug("hainan not show result in web,sleep..,", searchTime)
				sgthread.SleepBySecond(sleepTime)
				continue
			} else {
				sglog.Debug("hainan had show result in web,", searchTime)
			}
		}

		dataMap := make(map[string]string)

		searchPageStr := strconv.Itoa(searchPage)
		ignoreMap := make(map[string]string)
		totalPage, totalNum, err := dataSpider(title, index, searchTime, searchPageStr, ignoreMap)
		if err != nil {
			sglog.Error("search err,", searchTime, searchPageStr, err, "isPerson:", isPersonal)
			continue
		}
		if 0 == totalPage || 0 == totalNum {
			sglog.Error("not valid error,page or num is zero,page:", totalPage, ",num:", totalNum, "isPerson:", isPersonal)
			sgthread.SleepBySecond(60)
			continue
		}

		deleteRecord := 0
		//totalPage = 1
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
			//sglog.Debug("page:", i, "size:", len(tmpDataMap), ",totalsize:", len(dataMap))
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
			sglog.Error(title, "end search ", searchTime, ",needSize:", totalNum, ",real size:", len(dataMap), "isPerson:", isPersonal)
			sgthread.SleepBySecond(60)
			continue
		} else {
			sglog.Info(title, "end search ok", searchTime, ",needSize:", totalNum, ",real size:", len(dataMap), "isPerson:", isPersonal)
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

				if cache.IsCardDataExist(title, k) {
					sglog.Error("data already exist,time:", searchTime, "code:", k)
				} else {
					cardDataMap[k] = tmp
				}
			}
			cache.AddCardDataToMem(cardDataMap)

			downData := data.ChangeDownloadStatus(title, urlTips, define.DEF_DOWNLOAD_STATUS_COMPLETE, searchTime)
			db.UpdateDownloadToDb(downData)
			memberType := define.MEMBER_TYPE_COMPANY
			if isPersonal {
				memberType = define.MEMBER_TYPE_PERSIONAL
			}
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

				data.UpdateLastestInfo(title, define.CARD_TYPE_NORMAL, memberType, searchTime)
				checkNoticeToUser(title)
			}(cardDataMap)
		}
		startDt = startDt.AddDate(0, 1, 0)

		sgthread.SleepBySecond(5)
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
			if define.CITY_TIANJIN == title {
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
									sglog.Error("tranform error by tr1,", td.Text(), err)
								}
								totalNum, err = strconv.Atoi(strlist[1])
								if err != nil {
									sglog.Error("tranform error by tr2,", td.Text(), err)
								}
								//sglog.Info(td.Text(), "parse: page", totalPage, ",num:", totalNum)
							}
						}
					}
				})
			}
		}
	})

	if title == define.CITY_HAINAN {
		doc.Find("div,span,li,ul").Each(func(_ int, cl *goquery.Selection) {
			tmpS := cl.Text()
			tmpS = strings.Replace(tmpS, " ", "", -1)
			strlist := strings.Split(tmpS, "\n")
			//sglog.Debug("len:", len(strlist), strlist)
			if len(strlist) == 5 && strlist[1] != "" {
				//sglog.Debug("22:", cl.Text())
				totalNum, err = strconv.Atoi(strlist[1])
				if err != nil {
					sglog.Error("tranform error hainan1,", strlist[1], err, strlist)
				}
				totalPage, err = strconv.Atoi(strlist[2])
				if err != nil {
					sglog.Error("tranform error hainan2,", strlist[2], err, strlist)
				}
				//sglog.Info(strlist, "parse: page", totalPage, ",num:", totalNum)
			}
		})
	}

	return totalPage, totalNum, nil
}

func checkNoticeToUser(title string) {
	if isCurrentMonthAllUpdate(title) {
		nowTime := time.Now()
		curTimeStr := sgtime.YearString(&nowTime) + sgtime.MonthString(&nowTime)
		sglog.Info(title, " ex current month data all updates!!!!!")
		notice.NoticeCurrentMonthDataUpdate(title, curTimeStr)
	}
}

func isCurrentMonthAllUpdate(title string) bool {
	nowTime := time.Now()
	curTimeStr := sgtime.YearString(&nowTime) + sgtime.MonthString(&nowTime)
	curLastestInfo := data.GetLastestCardInfo(title)
	if curTimeStr == curLastestInfo.TimeStr {
		if curLastestInfo.IsAllCardInfoUpdate() {
			return true
		}
	}
	return false
}
