package spider

import (
	"net/http"
	"strings"
	"xcxYaohao/src/define"
	"xcxYaohao/src/httpHandle"

	"github.com/PuerkitoBio/goquery"
	"github.com/coderguang/GameEngine_go/sglog"
)

func TianjinOldDataSpider(cmd []string) {
	unitIndex := "http://apply.xkctk.jtys.tj.gov.cn/apply/norm/unitQuery.html"
	//persionIndex := "http://apply.xkctk.jtys.tj.gov.cn/apply/norm/personQuery.html"
	dataMap := make(map[string]string)

	totalPage, totalNum, err := dataSpider(unitIndex, "201907", "1", dataMap)

	if err != nil {
		return
	}
	sglog.Info("page:", totalPage, ",num:", totalNum, "data:", dataMap)

}

func dataSpider(index string, timestr string, page string, dataMap map[string]string) (string, string, error) {
	totalPage := ""
	totalNum := ""
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
					sglog.Info(code, ":", name)
					if httpHandle.CheckCodeValid(define.CITY_TIANJIN, code) {
						dataMap[code] = name
					}
				}
			})
		} else if tmpS.Size() == 3 {
			tr.Find("td").Each(func(ix int, td *goquery.Selection) {
				sglog.Info(ix, ":", td.Text())
			})
		}
	})
	return totalPage, totalPage, nil
}
