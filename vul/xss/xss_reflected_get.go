package xss

import (
	"net/http"
	"pikachu-go/templates"
)

func ReflectedGetHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		value := r.URL.Query().Get("keyword")
		data := templates.NewPageData2(7, 9, value) // 直接输出，模拟未过滤
		renderer.RenderPage(w, "xss/xss_reflected_get.html", data)
	}
}
