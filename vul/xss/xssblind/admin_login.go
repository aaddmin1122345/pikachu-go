package xssblind

import (
	"net/http"
	"pikachu-go/templates"
)

func AdminLoginHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewPageData2(7, 21, "")
		renderer.RenderPage(w, "admin_login.html", data)
	}
}
