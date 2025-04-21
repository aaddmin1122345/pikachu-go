package xssblind

import (
	"net/http"
	"pikachu-go/templates"
)

func XssBlindHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewPageData2(7, 19, "")
		renderer.RenderPage(w, "xss/xss_blind.html", data)
	}
}
