package main

import (
	"log"
	"xcxYaohao/src/config"
	"xcxYaohao/src/data"
	"xcxYaohao/src/dataMove"
	"xcxYaohao/src/db"
	"xcxYaohao/src/define"
	"xcxYaohao/src/httpHandle"
	"xcxYaohao/src/sms"
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
	sgcmd.RegistCmd("SendTestRandom", "[\"SendTestRandom\"] :send  random", sms.SendTestRandom)
	sgcmd.RegistCmd("SendTestLuck", "[\"SendTestLuck\"] :send  luck", sms.SendTestLuck)
	sgcmd.RegistCmd("SendTestUnLuck", "[\"SendTestUnLuck\"] :send  unluck", sms.SendTestUnLuck)
	sgcmd.RegistCmd("DataMove", "[\"DataMove\"] :move old data to new ", dataMove.InitDb)

}

func main() {

	sgcfg.SetServerCfgDir("./../globalConfig/xcxYaohao/")
	sgserver.StartServer(sgserver.ServerTypeLog, "debug", "./log/", log.LstdFlags, true)
	sgserver.StartServer(sgserver.ServerTypeMail)

	config.InitCfg()

	data.InitWxOpenIdCfg()

	spider.AutoCreateFileDir()

	spider.TianjinOldDataSpider([]string{})

	sgcmd.StartCmdWaitInputLoop()

	db.InitDb()

	titlelist := config.GetTitleList()
	for _, v := range titlelist {
		title := v
		if title != define.CITY_TIANJIN {
			continue
		}
		go spider.NewSpider(title)
	}

	go httpHandle.NewWebServer(config.GetUtilCfg().Port)

	go spider.InitClear()

	RegistCmd()
	sgcmd.StartCmdWaitInputLoop()
}
