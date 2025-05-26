package ssrf

import (
	"net/http"
	"pikachu-go/templates"
	"pikachu-go/vul/ssrf/ssrf_info"
)

// SsrfHandler 渲染 SSRF 漏洞概述页面
func SsrfHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewPageData2(105, 106, "")
		renderer.RenderPage(w, "ssrf/ssrf.html", data)
	}
}

// Info1Handler 返回第一个诗歌处理函数
func Info1Handler() http.HandlerFunc {
	return ssrf_info.Info1Handler()
}

// Info2Handler 返回第二个诗歌处理函数
func Info2Handler() http.HandlerFunc {
	return ssrf_info.Info2Handler()
}
