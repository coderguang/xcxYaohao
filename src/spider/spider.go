package spider

import (
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"xcxYaohao/src/config"
	"xcxYaohao/src/data"
	"xcxYaohao/src/db"
	"xcxYaohao/src/define"
	"xcxYaohao/src/httpHandle"

	"github.com/coderguang/GameEngine_go/sgfile"

	"github.com/coderguang/GameEngine_go/sgtime"

	"github.com/gocolly/colly/extensions"

	"github.com/gocolly/colly"

	"github.com/coderguang/GameEngine_go/sglog"
	"github.com/coderguang/GameEngine_go/sgregex"
	"github.com/coderguang/GameEngine_go/sgstring"
	"github.com/coderguang/GameEngine_go/sgthread"
)

const (
	VISTIS_TIME_DISTANCE      int = 200
	DOWNLOAD_TIME_DISTANCE    int = 500
	RE_VISTIS_TIME_DISTANCE   int = 200
	RE_DOWNLOAD_TIME_DISTANCE int = 500
)

const (
	PDF_FILE_DIR          string = "./data/pdf/"
	TXT_FILE_DIR          string = "./data/txt/"
	PDF_FILE_COMPLETE_DIR string = "./data/pdf_complete/"
	TXT_FILE_COMPLETE_DIR string = "./data/txt_complete/"
)

var (
	globalSpiderMap map[string]*Spider
)

func init() {
	globalSpiderMap = make(map[string]*Spider)
}

func NewSpider(title string) {

	if _, ok := globalSpiderMap[title]; ok {
		sglog.Error("spider already exist,title:", title)
		return
	}
	spider := new(Spider)
	globalSpiderMap[title] = spider
	spider.StartAutoVisitUrl(title)
}

func GetSpider(title string) (*Spider, error) {
	if v, ok := globalSpiderMap[title]; ok {
		return v, nil
	}
	return nil, errors.New("not this spider,title:" + title)
}

func AutoCreateFileDir() {
	titlelist := config.GetTitleList()
	for _, v := range titlelist {
		sgfile.AutoMkDir(PDF_FILE_DIR + v)
		sgfile.AutoMkDir(TXT_FILE_DIR + v)
		sgfile.AutoMkDir(PDF_FILE_COMPLETE_DIR + v)
		sgfile.AutoMkDir(TXT_FILE_COMPLETE_DIR + v)
	}
}

func (spider *Spider) StartAutoVisitUrl(title string) {
	cfg, err := config.GetSpiderCfg(title)
	if err != nil {
		sglog.Error("can't find ", title, " spider config")
		sgthread.DelayExit(2)
	}
	spider.cfg = cfg
	spider.hadVisitUrls = make(map[string]bool)
	spider.ignoreUrls = []string{}
	spider.ignoreUrls = append(spider.ignoreUrls, spider.cfg.IgnoreUrls...)
	spider.collector = colly.NewCollector()
	spider.collector.IgnoreRobotsTxt = true
	spider.collector.CheckHead = true
	spider.collector.AllowedDomains = spider.cfg.AllowUrls
	spider.collector.AllowURLRevisit = true
	spider.collector.WithTransport(&http.Transport{
		DisableKeepAlives: true,
	})

	extensions.RandomUserAgent(spider.collector)
	extensions.Referer(spider.collector)

	spider.collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		title := e.Text
		if !strings.Contains(link, spider.cfg.HTTP) {
			link = spider.cfg.HTTP + ":" + link
		}

		//sglog.Debug("find tilte:", title, ",link:", link)

		if !sgregex.URL(link) {
			//sglog.Debug(link, " not a valid url,title is", title)
			return
		}

		if spider.isIgnoreUrl(link) {
			return
		}
		if spider.isHadVisitUrl(link) {
			return
		}

		if sgstring.EqualWithOr(link, spider.ignoreUrls) {
			return
		}

		if !config.PageFliter(spider.cfg.Title, title) {
			if _, err = strconv.Atoi(title); err != nil {
				spider.hadVisitUrls[link] = true
				return
			}
		}

		if strings.Contains(title, "下一页") {
			spider.AddIgnoreUrl(link)
		}

		if strings.HasSuffix(link, "pdf") {
			sgthread.SleepByMillSecond(DOWNLOAD_TIME_DISTANCE)
			spider.DownloadFile(link, title)
		} else {
			sgthread.SleepByMillSecond(VISTIS_TIME_DISTANCE)
			er := e.Request.Visit(e.Request.AbsoluteURL(link))
			if er != nil {
				sglog.Error("start spider error by onHtml,url:=", e.Request.AbsoluteURL(link), ",err:=", er)
			}
		}
	})

	spider.collector.OnError(func(r *colly.Response, err error) {
		sglog.Error("visit ", r.Request.URL.String(), " occurt error, went wrong")
		sgthread.SleepByMillSecond(RE_VISTIS_TIME_DISTANCE)
		er := r.Request.Visit(r.Request.URL.String())
		if er != nil {
			sglog.Error("start spider error by onError,url:=", r.Request.URL.String(), ",err:=", er)
		}
		spider.hadVisitUrls[r.Request.URL.String()] = false
	})

	spider.collector.OnRequest(func(r *colly.Request) {
		//sglog.Info("Visiting ", r.URL.String(), " start...")
		spider.hadVisitUrls[r.URL.String()] = false
	})

	spider.collector.OnResponse(func(r *colly.Response) {
		sglog.Info("Visiting ", r.Request.URL.String(), " complete")
		spider.hadVisitUrls[r.Request.URL.String()] = true
	})

	sglog.Info("start loop spider ,title:", spider.cfg.Title)

	//spider.StartLoopSpider()
}

func (spider *Spider) StartLoopSpider() {
	sleepTime := 60
	for {
		//redownload
		downlist := data.GetReDownloadList(spider.cfg.Title)
		for _, v := range downlist {
			sglog.Info(spider.cfg.Title, " re download url:", v.URL)
			if !strings.Contains(v.Tips, "by reload") {
				v.Tips = "by reload" + v.Tips
			}
			spider.DownloadFile(v.URL, v.Tips)
			sgthread.SleepByMillSecond(RE_DOWNLOAD_TIME_DISTANCE)
		}

		//revisit
		revisitlist := spider.GetRevisitList()
		for _, v := range revisitlist {
			sglog.Info("re visist url:", v)
			spider.collector.Visit(v)
			sgthread.SleepByMillSecond(RE_VISTIS_TIME_DISTANCE)

		}

		er := spider.collector.Visit(spider.cfg.IndexURL)
		if er != nil {
			sglog.Error("start spider error,url:=", spider.cfg.IndexURL, ",err:=", er)
		}

		nowTime := time.Now()
		timeInt := time.Duration(300) * time.Second
		if 0 == len(downlist) && 0 == len(revisitlist) {

			normalTime := time.Date(nowTime.Year(), nowTime.Month(), 26, 9, 0, 0, 0, nowTime.Location())

			if nowTime.Before(normalTime) {
				timeInt = normalTime.Sub(nowTime)
			} else {

				// check is current month day all get
				curLastestInfo := data.GetLastestCardInfo(spider.cfg.Title)
				curTimeStr := sgtime.YearString(&nowTime) + sgtime.MonthString(&nowTime)

				curMonthAllUpdate := false
				if curTimeStr == curLastestInfo.TimeStr {
					if curLastestInfo.IsAllCardInfoUpdate() {
						curMonthAllUpdate = true
					}
				}

				if curMonthAllUpdate {
					//
					sglog.Info("current month data all updates!!!!!")
					httpHandle.NoticeCurrentMonthDataUpdate(spider.cfg.Title, curTimeStr)
					nextMonthDt := normalTime.AddDate(0, 1, 0)
					timeInt = nextMonthDt.Sub(nowTime)
				} else {

					hour := time.Now().Hour()
					if hour < 9 {
						nextRun := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 9, 0, 0, 0, nowTime.Location())
						timeInt = nextRun.Sub(nowTime)
					} else if hour > 19 {
						nextRun := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 23, 59, 59, 0, nowTime.Location())
						timeInt = nextRun.Sub(nowTime)
					}
				}
			}
			sleepTime = int(timeInt/time.Second) + 1
		}
		sglog.Info(spider.cfg.Title, " data collection now in sleep,will run after ", sleepTime, "s,", nowTime.Add(timeInt))
		sgthread.SleepBySecond(sleepTime)
	}
}

func (spider *Spider) DownloadFile(url string, title string) error {
	if !data.NeedDownloadFile(spider.cfg.Title, url) {
		return nil
	}
	downloadflag := false

	defer func() {
		if !downloadflag {
			downData := data.ChangeDownloadStatus(spider.cfg.Title, url, define.DEF_DOWNLOAD_STATUS_ERROR, title)
			db.UpdateDownloadToDb(downData)
		}
	}()

	sglog.Info("start download pdf,title:", spider.cfg.Title, ",url:", url)

	downData := data.ChangeDownloadStatus(spider.cfg.Title, url, define.DEF_DOWNLOAD_STATUS_DOWNING, title)
	db.UpdateDownloadToDb(downData)

	rawFileName, pdfFileName, err := spider.DownloadFileFromWeb(url, title)
	if err != nil {
		return err
	}
	err = spider.TransportPDFToTxt(rawFileName, pdfFileName)
	if err != nil {
		return err
	}

	txtFileName := strings.Replace(rawFileName, "pdf", "txt", -1)
	txtFileName = TXT_FILE_DIR + spider.cfg.Title + "/" + txtFileName

	timestr, memberType, cardType, updateNum, err := spider.ReadTxtFileAndInsertToDb(txtFileName)
	if err != nil {
		sglog.Error("read txt file error,", txtFileName, err)
		return err
	}

	sglog.Info("download file ok,", title, url, timestr, memberType, cardType, ",update num:", updateNum)

	downData = data.ChangeDownloadStatus(spider.cfg.Title, url, define.DEF_DOWNLOAD_STATUS_COMPLETE, title)
	db.UpdateDownloadToDb(downData)

	spider.RemoveAndRenameFile(pdfFileName, txtFileName, timestr, memberType, cardType)

	downloadflag = true
	return nil
}

func (spider *Spider) DownloadFileFromWeb(url string, title string) (string, string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	res, err := client.Get(url)
	if err != nil {
		sglog.Error("download file error,all try failed,title:", spider.cfg.Title, "src:", url, " ,error", err)
		if res != nil {
			res.Body.Close()
		}
		return "", "", err
	}
	defer res.Body.Close()
	str := strings.Split(url, "/")

	if len(str) <= 0 {
		return "", "", err
	}

	rawFileName := str[len(str)-1]
	pdfFileName := PDF_FILE_DIR + spider.cfg.Title + "/" + rawFileName

	f, err := os.Create(pdfFileName)
	if err != err {
		sglog.Error("create file error,title:", spider.cfg.Title, ",src:", url, ",error:", err)
		return "", "", err
	}
	io.Copy(f, res.Body)
	sglog.Info("download file ", url, " ,complete,file in ", pdfFileName)

	return rawFileName, pdfFileName, nil
}

func (spider *Spider) TransportPDFToTxt(rawFileName string, pdfFileName string) error {
	sglog.Info("start transform pdf to txt")

	cmd := exec.Command("python3", "index.py", spider.cfg.Title, rawFileName)
	sglog.Info("command is python3 index.py", spider.cfg.Title, " ", rawFileName)
	out, err := cmd.Output()
	if err != nil {
		sglog.Error("exec parse pdf to txt by python error,file=", pdfFileName, ",err=", err)
		return err
	}
	sglog.Info("output is :\n", string(out))
	return nil
}

func (spider *Spider) ReadTxtFileAndInsertToDb(fileDir string) (string, int, int, int, error) {
	contents, err := sgfile.GetFileContentAsStringLines(fileDir)
	if err != nil {
		sglog.Error("read txt file error,file:", fileDir, ",err=", err)
		return "", 0, 0, 0, err
	}
	//datas, timestr, totalnum, cardType, memberType, err
	datas, timestr, totalnum, cardType, memberType, err := config.TxtFileFliter(spider.cfg.Title, contents)

	if err != nil {
		sglog.Error("parse txt filter error,title:", spider.cfg.Title, ",file:", fileDir, err)
		return "", 0, 0, 0, err
	}

	sglog.Info("parse txt file detail complete,title:", spider.cfg.Title, ",timestr:", timestr, ",num:", totalnum, ",cardType:", cardType, ",memberType:", memberType)
	updateNum := 0
	updateMap := make(map[string]*define.CardData)
	for _, v := range datas {
		if !data.IsDataExist(v.Title, v.Code) {
			updateNum++
			updateMap[v.Code] = v
		}
	}

	now := sgtime.New()
	data.AddCardData(updateMap)
	data.UpdateLastestInfo(spider.cfg.Title, cardType, memberType, timestr)
	end := sgtime.New()
	sglog.Info("update card data in member size:", updateNum, ",use time:", (sgtime.GetTotalSecond(end) - sgtime.GetTotalSecond(now)))

	go func(updates map[string]*define.CardData) {
		nowEx := sgtime.New()
		updateDbNum := 0
		for _, v := range updates {
			if err = db.UpdateCardData(v); err != nil {
				sglog.Error("update data to db error,title", v.Title, "code", v.Code, err)
			} else {
				updateDbNum++
			}
		}
		endEx := sgtime.New()
		sglog.Info(spider.cfg.Title, " update card data in databases size:", updateDbNum, ",use time:", (sgtime.GetTotalSecond(endEx) - sgtime.GetTotalSecond(nowEx)))

	}(updateMap)

	return timestr, memberType, cardType, updateNum, nil
}

func (spider *Spider) RemoveAndRenameFile(pdfFileName string, txtFileName string, timestr string, memberType int, cardType int) {

	memberTypeStr := "persional"
	if memberType == define.MEMBER_TYPE_COMPANY {
		memberTypeStr = "company"
	}
	cardTypeStr := "normal"
	if cardType == define.CARD_TYPE_COMPANY {
		cardTypeStr = "conservation"
	}
	rename := timestr + "_" + memberTypeStr + "_" + cardTypeStr
	completePDFDir := PDF_FILE_COMPLETE_DIR + spider.cfg.Title + "/" + timestr + "/"
	sgfile.AutoMkDir(completePDFDir)
	newPDFFile := completePDFDir + rename + ".pdf"
	err := sgfile.Rename(pdfFileName, newPDFFile)
	if err != nil {
		sglog.Error("remove file error,", pdfFileName, "=========>", newPDFFile, err)
	}

	completeTxtDir := TXT_FILE_COMPLETE_DIR + spider.cfg.Title + "/" + timestr + "/"
	sgfile.AutoMkDir(completeTxtDir)
	newTxtFile := completeTxtDir + rename + ".txt"
	err = sgfile.Rename(txtFileName, newTxtFile)
	if err != nil {
		sglog.Error("remove file error,", txtFileName, "=========>", newTxtFile, err)
	}

}
