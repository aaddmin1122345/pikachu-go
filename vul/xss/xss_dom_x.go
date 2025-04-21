package xss

import (
	"net/http"
	"pikachu-go/templates"
)

func DomXssXHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewPageData2(7, 16, "")
		renderer.RenderPage(w, "xss/xss_dom_x.html", data)
	}
}
