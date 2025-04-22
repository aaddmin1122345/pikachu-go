package xss

import (
	"fmt"
	"net/http"
	"pikachu-go/templates"
)

func ReflectedGetHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		htmlMsg := ""

		if r.URL.Query().Get("submit") != "" {
			message := r.URL.Query().Get("message")
			if message == "" {
				htmlMsg = "<p class='notice'>输入'kobe'试试-_-</p>"
			} else if message == "kobe" {
				htmlMsg = fmt.Sprintf("<p class='notice'>愿你和%s一样，永远年轻，永远热血沸腾！</p><img src='/assets/images/nbaplayer/kobe.png' />", message)
			} else {
				htmlMsg = fmt.Sprintf("<p class='notice'>who is %s, i don't care!</p>", message)
			}
		}

		data := templates.NewPageData2(7, 9, htmlMsg)
		renderer.RenderPage(w, "xss/xss_reflected_get.html", data)
	}
}
