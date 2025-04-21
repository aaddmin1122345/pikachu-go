package fileinclude

import (
	"net/http"
	"pikachu-go/templates"
)

// FileIncludeHandler 渲染 File Inclusion 漏洞概述页面
func FileIncludeHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewPageData2(74, 75, "")
		renderer.RenderPage(w, "fileinclude/fileinclude.html", data)
	}
}
