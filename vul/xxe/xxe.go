package xxe

import (
	"net/http"
	"pikachu-go/templates"
)

// XxeHandler 渲染 XXE 概述页面（对应 xxe.php）
func XxeHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewPageData2(77, 78, "")
		renderer.RenderPage(w, "xxe/xxe.html", data)
	}
}
