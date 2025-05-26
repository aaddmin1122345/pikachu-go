package xss

import (
	"html/template"
	"net/http"
	"pikachu-go/templates"
)

func DomXssXHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewPageData2(7, 13, "")

		// 如果有text参数，添加一个点击链接
		if r.URL.Query().Get("text") != "" {
			data.HtmlMsg = template.HTML("<a href='#' onclick='domxss()'>有些费尽心机想要忘记的事情,后来真的就忘掉了</a>")
		}

		renderer.RenderPage(w, "xss/xss_dom_x.html", data)
	}
}
