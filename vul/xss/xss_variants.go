package xss

import (
	"html"
	"net/http"
	"pikachu-go/templates"
	"strings"
)

func RenderXssVariant(renderer templates.Renderer, variant string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input := r.URL.Query().Get("keyword")
		var output string

		subIndex := 31 // 默认 xss_01
		switch variant {
		case "xss_01":
			subIndex = 31
			output = input
		case "xss_02":
			subIndex = 32
			output = html.EscapeString(input)
		case "xss_03":
			subIndex = 33
			output = strings.ReplaceAll(strings.ReplaceAll(input, "<", "&lt;"), ">", "&gt;")
		case "xss_04":
			subIndex = 34
			output = html.EscapeString(input)
		}
		data := templates.NewPageData2(7, subIndex, output)
		renderer.RenderPage(w, "xss/"+variant+".html", data)

	}
}
