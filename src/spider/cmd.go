package spider

import (
	"strconv"
	"xcxYaohao/src/data"
	"xcxYaohao/src/db"

	"github.com/coderguang/GameEngine_go/sglog"
)

func ShowIgnors(cmd []string) {
	title := cmd[1]
	spider, err := GetSpider(title)
	if err != nil {
		sglog.Error("no this spider,show ignores error", title, err)
		return
	}
	spider.ShowIgnoreUrls()
}

func TransportPDFToTxt(cmd []string) {
	title := cmd[1]
	spider, err := GetSpider(title)
	if err != nil {
		sglog.Error("no this spider,TransportPDFToTxt error", title, err)
		return
	}
	rawFileName := cmd[2]
	pdfFileName := PDF_FILE_DIR + spider.cfg.Title + "/" + rawFileName
	err = spider.TransportPDFToTxt(rawFileName, pdfFileName)
	if err != nil {
		sglog.Error("TransportPDFToTxt error", title, rawFileName, err)
		return
	}
}

func ReadTxtFileToDb(cmd []string) {
	title := cmd[1]
	spider, err := GetSpider(title)
	if err != nil {
		sglog.Error("no this spider,ReadTxtFileToDb error", title, err)
		return
	}
	rawFileName := cmd[2]
	txtFileName := TXT_FILE_DIR + spider.cfg.Title + "/" + rawFileName

	timestr, memberType, cardType, updateNum, err := spider.ReadTxtFileAndInsertToDb(txtFileName)
	if err != nil {
		sglog.Error("ReadTxtFileToDb read txt file error,", rawFileName, err)
		return
	}

	sglog.Info("ReadTxtFileToDb file ok,", title, timestr, memberType, cardType, ",update num:", updateNum)

}

func ChangeDownloadStatus(cmd []string) {
	title := cmd[1]
	spider, err := GetSpider(title)
	if err != nil {
		sglog.Error("no this spider,ChangeDownloadStatus error", title, err)
		return
	}
	url := cmd[2]
	status := cmd[3]

	statusInt, err := strconv.Atoi(status)
	if err != nil {
		sglog.Error("ChangeDownloadStatus error", title, "status error,status is ", status, err)
		return
	}
	downData := data.ChangeDownloadStatus(spider.cfg.Title, url, statusInt, "change by cmd")
	db.UpdateDownloadToDb(downData)

}
