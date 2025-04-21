package ssrf

import (
	"net/http"
	"pikachu-go/templates"
)

// SsrfHandler 渲染 SSRF 漏洞概述页面
func SsrfHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewPageData2(60, 61, "")
		renderer.RenderPage(w, "ssrf/ssrf.html", data)
	}
}
