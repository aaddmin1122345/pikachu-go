package csrf

import (
	"net/http"
	"pikachu-go/templates"
)

// CsrfHandler 渲染 CSRF 漏洞模块首页
func CsrfHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewPageData2(25, 26, "")
		renderer.RenderPage(w, "csrf/csrf.html", data)
	}
}
