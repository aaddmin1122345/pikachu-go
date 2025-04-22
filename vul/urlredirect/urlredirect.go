package urlredirect

import (
	"net/http"
	"pikachu-go/templates"
)

func URLRedirectHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		html := ""

		if url != "" {
			if url == "i" {
				html += "<p>好的,希望你能坚持做你自己!</p>"
			} else {
				http.Redirect(w, r, url, http.StatusFound)
				return
			}
		}

		data := templates.NewPageData2(100, 102, html) // index 102 for urlredirect
		renderer.RenderPage(w, "urlredirect/urlredirect.html", data)
	}
}
