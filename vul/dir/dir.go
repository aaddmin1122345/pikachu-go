package dir

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"pikachu-go/templates"
)

func DirHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("dir")
		if query == "" {
			query = "."
		}
		query = filepath.Clean(query)

		html := "<ul>"
		files, err := os.ReadDir(query)
		if err == nil {
			for _, file := range files {
				name := file.Name()
				displayPath := strings.TrimPrefix(filepath.Join(query, name), "./")

				if file.IsDir() {
					html += `<li>[DIR] <a href="?dir=` + displayPath + `">` + name + `</a></li>`
				} else {
					html += `<li>[FILE] <a href="/vul/dir/dir_list?title=` + displayPath + `">` + name + `</a></li>`
				}
			}
		}
		html += "</ul>"

		data := templates.NewPageData2(80, 81, html)
		renderer.RenderPage(w, "dir/dir.html", data)
	}
}
