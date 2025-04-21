package dir

import (
	"net/http"
	"pikachu-go/templates"
)

// DirHandler 渲染目录遍历漏洞概述页面（仅说明性展示）
func DirHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewPageData2(80, 81, "")
		renderer.RenderPage(w, "dir/dir.html", data)
	}
}
