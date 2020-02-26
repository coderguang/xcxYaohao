package main

import (
	"log"
	"xcxYaohao/src/cache"
	"xcxYaohao/src/config"
	"xcxYaohao/src/data"
	"xcxYaohao/src/dataMove"
	"xcxYaohao/src/db"
	"xcxYaohao/src/httpHandle"
	"xcxYaohao/src/notice"
	"xcxYaohao/src/sms"
	"xcxYaohao/src/spider"

	"github.com/coderguang/GameEngine_go/sgcfg"
	"github.com/coderguang/GameEngine_go/sgcmd"
	"github.com/coderguang/GameEngine_go/sgmail"
	"github.com/coderguang/GameEngine_go/sgserver"
)

func TestSendMail(cmd []string) {
	sgmail.SendMail("xcxYaohao test", []string{config.GetUtilCfg().Receiver}, "test mail")
}

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
	sgcmd.RegistCmd("ShowCurrentDatas", "[\"ShowCurrentDatas\"]:show current statistic data", data.ShowCurrentDatas)
	sgcmd.RegistCmd("ReloadBoardcast", "[\"ReloadBoardcast\"] :ReloadBoardcast cfg", data.ReloadBoardcast)
	sgcmd.RegistCmd("NoticeSmsByCmd", "[\"NoticeSmsByCmd\",\"shenzhen\",\"201908\"] :notice sms", notice.NoticeSmsByCmd)
	sgcmd.RegistCmd("ShowCurrentSmsFlang", "[\"ShowCurrentSmsFlang\"] :ShowCurrentSmsFlang", sms.ShowCurrentSmsFlang)
	sgcmd.RegistCmd("changeCurrentSmsFlang", "[\"changeCurrentSmsFlang\"] :changeCurrentSmsFlang ", sms.ChangeCurrentSmsFlang)
	sgcmd.RegistCmd("TestSendMail", "[\"TestSendMail\"] :TestSendMail ", TestSendMail)
	sgcmd.RegistCmd("TianJinOldData", "[\"TianJinOldData\"] :TianJinOldData spider data", spider.TianjinOldDataSpider)
}

func main() {

	sgcfg.SetServerCfgDir("./../globalConfig/xcxYaohao/")
	sgserver.StartServer(sgserver.ServerTypeLog, "debug", "./log/", log.LstdFlags, true)
	sgserver.StartServer(sgserver.ServerTypeMail)

	config.InitCfg()
	data.ReloadBoardcast([]string{})

	data.InitOpenIdCfgs()

	spider.AutoCreateFileDir()

	db.InitDb()

	cache.InitAndLoadCardData()

	//hainan special logic
	go spider.HainanOldDataSpider([]string{})

	titlelist := config.GetTitleList()
	for _, v := range titlelist {
		title := v
		go spider.NewSpider(title)
	}

	go httpHandle.NewWebServer(config.GetUtilCfg().Port)

	go spider.InitClear()

	RegistCmd()
	sgcmd.StartCmdWaitInputLoop()
}
