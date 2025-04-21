package unsafedownload

import (
	"net/http"
	"pikachu-go/templates"
)

// UnsafedownloadHandler 对应 unsafedownload.php，渲染漏洞概述页
func UnsafedownloadHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewPageData2(52, 53, "")
		renderer.RenderPage(w, "unsafedownload/unsafedownload.html", data)
	}
}
