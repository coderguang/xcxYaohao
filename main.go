package main

import (
	"log"
	"xcxYaohao/src/config"
	"xcxYaohao/src/db"
	"xcxYaohao/src/spider"

	"github.com/coderguang/GameEngine_go/sgcfg"
	"github.com/coderguang/GameEngine_go/sgcmd"
	"github.com/coderguang/GameEngine_go/sgserver"
)

func main() {

	sgcfg.SetServerCfgDir("./../globalConfig/xcxYaohao/")
	sgserver.StartServer(sgserver.ServerTypeLog, "debug", "./../log/", log.LstdFlags, true)
	sgserver.StartServer(sgserver.ServerTypeMail)

	config.InitCfg()

	db.InitDb()

	spider.AutoCreateFileDir()

	titlelist := config.GetTitleList()
	spiderlist := []spider.Spider{}
	for _, v := range titlelist {
		tmpSpider := spider.Spider{}
		tmpSpider.StartAutoVisitUrl(v)
		spiderlist = append(spiderlist, tmpSpider)
	}

	sgcmd.StartCmdWaitInputLoop()

}
