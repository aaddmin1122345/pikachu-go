package unsafedownload

import (
	"net/http"
	"pikachu-go/templates"
)

// DownNbaHandler 对应 down_nba.php，展示球星头像下载入口
func DownNbaHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewPageData2(52, 54, "")
		renderer.RenderPage(w, "unsafedownload/down_nba.html", data)
	}
}
