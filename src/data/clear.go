package data

import (
	"xcxYaohao/src/define"

	"github.com/coderguang/GameEngine_go/sglog"
	"github.com/coderguang/GameEngine_go/sgtime"
	"github.com/mohae/deepcopy"
)

func ClearData() *define.StatisticsData {
	globalOpenIds.Lock.Lock()
	now := sgtime.New()
	idset := make(map[string]bool)
	for k, v := range globalOpenIds.Data {
		if _, ok := idset[v.Openid]; !ok {
			idset[v.Openid] = true
		}
		if sgtime.GetTotalSecond(now)-sgtime.GetTotalSecond(sgtime.TransfromTimeToDateTime(*v.Time)) > 3600 {
			sglog.Debug("delete openid:", v.Openid, ",code:", v.Code)
			delete(globalOpenIds.Data, k)
		}
	}
	globalOpenIds.Lock.Unlock()
	sglog.Info("clear openid data complete,user size:", len(idset), globalStatics)

	//sgmail.SendMail("xcxYaohao statistic", []string{config.GetUtilCfg().Receiver}, globalStatics.String())

	tmp := deepcopy.Copy(globalStatics)

	globalStatics.Reset()

	tmpV, ok := tmp.(*define.StatisticsData)
	if ok {
		return tmpV
	}
	return nil
}