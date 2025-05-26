package infoleak

import (
	"net/http"
	"pikachu-go/templates"
)

// InfoleakHandler 渲染敏感信息泄露模块介绍页面
func InfoleakHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewPageData2(85, 86, "")
		renderer.RenderPage(w, "infoleak/infoleak.html", data)
	}
}
