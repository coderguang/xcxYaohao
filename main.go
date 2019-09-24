package main

import (
	"log"
	"xcxYaohao/src/config"
	"xcxYaohao/src/db"

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

	sgcmd.StartCmdWaitInputLoop()

}
