package infoleak

import (
	"net/http"
	"pikachu-go/templates"
)

// InfoleakHandler 渲染敏感信息泄露模块介绍页面
func InfoleakHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewPageData2(66, 67, "") // 激活左侧导航栏
		renderer.RenderPage(w, "infoleak/infoleak.html", data)
	}
}
