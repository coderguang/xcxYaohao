package config

import (
	"strings"

	"github.com/coderguang/GameEngine_go/sgstring"
)

//筛选一个页面是否应该被访问
func PageFliter(title, pageTitle string) bool {
	switch title {
	case "guangzhou":
		return guangzhouFliter(pageTitle)
	case "shenzhen":
		return shenzhenFliter(pageTitle)
	}
	return false
}

func guangzhouFliter(pageTitle string) bool {
	pageTxt := []string{"配置结果"}
	if !sgstring.ContainsWithOr(pageTitle, pageTxt) && !strings.Contains(pageTitle, "下一页") && !strings.Contains(pageTitle, "pdf") {
		return false
	}
	return true
}

func shenzhenFliter(pageTitle string) bool {
	pageTxt := []string{"增量指标摇号结果公告", "指标配置结果"}
	if !sgstring.ContainsWithOr(pageTitle, pageTxt) && !strings.Contains(pageTitle, "下一页") && !strings.Contains(pageTitle, "pdf") {
		return false
	}
	return true
}
