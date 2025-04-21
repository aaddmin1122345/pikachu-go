package dir

import (
	"net/http"
	"pikachu-go/templates"
)

// DirListHandler 模拟页面列表跳转
func DirListHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Query().Get("title")
		data := templates.NewPageData2(4, 19, "")

		switch title {
		case "jarheads":
			renderer.RenderPage(w, "dir/soup/jarheads.html", data)
		case "truman":
			renderer.RenderPage(w, "dir/soup/truman.html", data)
		default:
			renderer.RenderPage(w, "dir/dir_list.html", data)
		}
	}
}
