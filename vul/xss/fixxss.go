package xss

import (
	"html/template"
	"net/http"
	"pikachu-go/templates"
)

func FixXssHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		output := ""
		if r.Method == http.MethodPost {
			input := r.FormValue("msg")
			output = template.HTMLEscapeString(input)
		}

		data := templates.NewPageData2(7, 17, output)
		renderer.RenderPage(w, "xss/fixxss.html", data)
	}
}
