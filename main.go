package main

import (
	"log"
	"xcxYaohao/src/config"
	"xcxYaohao/src/data"
	"xcxYaohao/src/db"
	"xcxYaohao/src/httpHandle"
	"xcxYaohao/src/spider"

	"github.com/coderguang/GameEngine_go/sgcfg"
	"github.com/coderguang/GameEngine_go/sgcmd"
	"github.com/coderguang/GameEngine_go/sgserver"
)

func RegistCmd() {
	sgcmd.RegistCmd("ShowIgnors", "[\"ShowIgnors\",\"guangzhou\"] :show ignores", spider.ShowIgnors)
	sgcmd.RegistCmd("PDF", "[\"PDF\",\"guangzhou\",\"1459152795388.pdf\"]:transform a pdf file to txt file", spider.TransportPDFToTxt)
	sgcmd.RegistCmd("TXT", "[\"TXT\",\"guangzhou\",\"1459152795388.txt\"]:get txt file and insert it to db", spider.ReadTxtFileToDb)
	sgcmd.RegistCmd("DownloadStatus", "[\"DownloadStatus\",\"shenzhen\",\"http://xqctk.jtys.sz.gov.cn/attachment/20160328/1459152795388.pdf\",\"1\",]:change download status", spider.ChangeDownloadStatus)
	sgcmd.RegistCmd("ShowLasteTime", "[\"ShowLasteTime\"] :show current", data.ShowLastestInfo)
}

func main() {

	sgcfg.SetServerCfgDir("./../globalConfig/xcxYaohao/")
	sgserver.StartServer(sgserver.ServerTypeLog, "debug", "./log/", log.LstdFlags, true)
	sgserver.StartServer(sgserver.ServerTypeMail)

	config.InitCfg()
	data.InitWxOpenIdCfg()
	db.InitDb()

	spider.AutoCreateFileDir()

	titlelist := config.GetTitleList()
	for _, v := range titlelist {
		title := v
		go spider.NewSpider(title)
	}

	go httpHandle.NewWebServer(config.GetUtilCfg().Port)

	go data.InitClear()

	RegistCmd()
	sgcmd.StartCmdWaitInputLoop()
}
