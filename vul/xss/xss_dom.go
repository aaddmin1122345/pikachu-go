package xss

import (
	"net/http"
	"pikachu-go/templates"
)

// DomXssHandler 处理DOM XSS请求
func DomXssHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewPageData2(7, 12, "")
		// 这里的GET参数实际上在后端不会用到，DOM XSS是在前端发生的
		if r.URL.Query().Get("text") != "" {
			// 在实际场景中，这里可能会有一些后端处理，但在演示中我们不做处理
			// 这是为了表明DOM XSS不依赖于后端输出
		}

		renderer.RenderPage(w, "xss/xss_dom.html", data)
	}
}
