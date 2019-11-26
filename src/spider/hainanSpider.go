package spider

import (
	"net/http"
	"strings"
	"xcxYaohao/src/define"

	"github.com/coderguang/GameEngine_go/sgstring"

	"github.com/PuerkitoBio/goquery"
	"github.com/coderguang/GameEngine_go/sgtime"

	"github.com/coderguang/GameEngine_go/sglog"
)

func HainanOldDataSpider(cmd []string) {
	unitIndex := "https://apply.hnjdctk.gov.cn/apply/app/status/norm/unit"
	persionIndex := "https://apply.hnjdctk.gov.cn/apply/app/status/norm/person"

	startDt, err := sgtime.ParseInLocation(sgtime.FORMAT_TIME_NORMAL, "2018-08-01 00:00:00")
	if err != nil {
		sglog.Error("parse startDt error,", err)
	}

	go StartSpiderEx(define.CITY_HAINAN, startDt, unitIndex, false)
	go StartSpiderEx(define.CITY_HAINAN, startDt, persionIndex, true)

}

func HainanIsShowResultInWebIndex(searchTimeStr string) bool {
	index := "https://www.hnjdctk.gov.cn/tzgg/"

	params := ""
	resp, err := http.Post(index,
		"application/x-www-form-urlencoded",
		strings.NewReader(params))

	if err != nil {
		sglog.Error("hainan check web index is show result,err:", err)
		return false
	}

	doc, err := goquery.NewDocumentFromResponse(resp)

	if err != nil {
		sglog.Error("hainan check web index is show result docs error")
		return false
	}
	isShowResult := false
	keyStr := []string{"海南省第", "期小客车增量指标配置结果公告"}
	doc.Find("table tr").Each(func(_ int, tr *goquery.Selection) {
		tr.Find("dd").Each(func(ix int, td *goquery.Selection) {
			title := td.Text()
			if sgstring.ContainsWithAnd(title, keyStr) {
				//sglog.Debug("title:", title)
				title = strings.Replace(title, "\t", "", -1)
				strlist := strings.Split(title, "\n")
				newTitleList := []string{}
				for _, v := range strlist {
					if v != "" {
						newTitleList = append(newTitleList, v)
					}
				}
				if len(newTitleList) == 2 {
					resultTime := newTitleList[1]
					timelist := strings.Split(resultTime, "-")
					if len(timelist) == 3 {
						timestr := timelist[0] + timelist[1]
						if timestr == searchTimeStr {
							isShowResult = true
							sglog.Info("hainan had show result in index,time=", timestr)
						}
					}
				}
			}
		})

	})

	return isShowResult
}
