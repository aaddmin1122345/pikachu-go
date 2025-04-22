package overpermission

import (
	"net/http"
	"pikachu-go/templates"
)

// OpHandler 渲染越权漏洞概述页面
func OpHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewPageData2(73, 74, "")
		renderer.RenderPage(w, "overpermission/op.html", data)
	}
}
