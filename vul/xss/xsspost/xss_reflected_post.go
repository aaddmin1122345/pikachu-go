package xsspost

import (
	"net/http"
	"pikachu-go/templates"
)

func XssReflectedPostHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := ""
		if r.Method == http.MethodPost {
			msg = r.FormValue("message")
		}

		data := templates.NewPageData2(7, 18, msg)
		renderer.RenderPage(w, "xss/xsspost/xss_reflected_post.html", data)
	}
}
