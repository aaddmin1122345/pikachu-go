package rce

import (
	"net/http"
	"pikachu-go/templates"
)

func RceIndexHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewPageData2(50, 51, "")
		renderer.RenderPage(w, "rce/rce.html", data)
	}
}
