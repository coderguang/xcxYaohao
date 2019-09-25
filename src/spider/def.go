package spider

import (
	"xcxYaohao/src/config"

	"github.com/coderguang/GameEngine_go/sglog"

	"github.com/gocolly/colly"
)

type Spider struct {
	cfg          config.SpiderCfg
	collector    *colly.Collector
	hadVisitUrls map[string]bool
	ignoreUrls   []string
}

func (data *Spider) ShowIgnoreUrls() {
	sglog.Info("------ignore---", data.cfg.Title, "---------")
	for _, v := range data.ignoreUrls {
		sglog.Debug(v)
	}
	sglog.Info("------ignore---", data.cfg.Title, "--complete-------")
}

func (data *Spider) AddIgnoreUrl(url string) {
	if data.isIgnoreUrl(url) {
		return
	}
	data.ignoreUrls = append(data.ignoreUrls, url)
}

func (data *Spider) isIgnoreUrl(url string) bool {
	for _, v := range data.ignoreUrls {
		if v == url {
			return true
		}
	}
	return false
}

func (data *Spider) isHadVisitUrl(url string) bool {
	if v, ok := data.hadVisitUrls[url]; ok {
		return v
	}
	return false
}

func (data *Spider) GetRevisitList() []string {
	relist := []string{}
	for k, v := range data.hadVisitUrls {
		if !v {
			relist = append(relist, k)
		}
	}
	return relist
}
