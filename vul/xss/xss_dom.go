package xss

import (
	"net/http"
	"pikachu-go/templates"
)

func DomXssHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewPageData2(7, 15, "")
		renderer.RenderPage(w, "xss/xss_dom.html", data)
	}
}
