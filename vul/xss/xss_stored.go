package xss

import (
	"net/http"
	"pikachu-go/templates"
)

var storedMsgs []string

func StoredHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			msg := r.FormValue("msg")
			storedMsgs = append(storedMsgs, msg)
		}

		// 拼接所有已存内容，模拟留言墙
		html := ""
		for _, m := range storedMsgs {
			html += m + "<br>"
		}

		data := templates.NewPageData2(7, 10, html)
		renderer.RenderPage(w, "xss/xss_stored.html", data)
	}
}
