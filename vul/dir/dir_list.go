package dir

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"pikachu-go/templates"
)

func DirListHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Query().Get("title")
		if title == "" {
			title = "index.php"
		}
		title = filepath.Clean(title)

		path := filepath.Join("vul/dir/soup", title)
		content := ""

		if data, err := os.ReadFile(path); err == nil {
			content = "<pre>" + template.HTMLEscapeString(string(data)) + "</pre>"
		}

		data := templates.NewPageData2(80, 82, content)
		renderer.RenderPage(w, "dir/dir_list.html", data)
	}
}
