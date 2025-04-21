package sqli

import (
	"net/http"
	"pikachu-go/templates"
)

// SqliHandler 渲染 SQL 注入漏洞概述页面
func SqliHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewPageData2(3, 4, "")
		renderer.RenderPage(w, "sqli/sqli.html", data)
	}
}
