package xss

import (
	"fmt"
	"html/template"
	"net/http"
	"pikachu-go/templates"
)

// ReflectedGetHandler 处理反射型XSS（GET）请求
func ReflectedGetHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewPageData2(7, 9, "")

		if r.URL.Query().Get("submit") != "" {
			message := r.URL.Query().Get("message")
			if message == "" {
				data.HtmlMsg = template.HTML("<p class='notice'>输入'kobe'试试-_-</p>")
			} else if message == "kobe" {
				data.HtmlMsg = template.HTML(fmt.Sprintf("<p class='notice'>愿你和%s一样，永远年轻，永远热血沸腾！</p><img src='/assets/images/nbaplayer/kobe.png' />", message))
			} else {
				// XSS漏洞点：直接输出用户输入而不进行过滤
				data.HtmlMsg = template.HTML(fmt.Sprintf("<p class='notice'>who is %s, i don't care!</p>", message))
			}
		}

		renderer.RenderPage(w, "xss/xss_reflected_get.html", data)
	}
}
