package xss

import (
	"net/http"
	"pikachu-go/templates"
)

// XssIndexHandler 渲染 xss 模块首页
func XssIndexHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewPageData2(7, 8, "")
		renderer.RenderPage(w, "xss/xss.html", data)
	}
}
