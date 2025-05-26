package urlredirect

import (
	"net/http"
	"pikachu-go/templates"
	"strings"
)

// URLRedirectHandler 处理URL重定向功能
func URLRedirectHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		html := ""

		if url != "" {
			if url == "i" {
				html = "<p class='notice'>好的,希望你能坚持做你自己!</p>"
			} else {
				// 确保URL以斜杠开头或包含完整协议
				if !strings.HasPrefix(url, "/") && !strings.Contains(url, "://") {
					url = "/" + url
				}

				// 重定向到指定URL
				http.Redirect(w, r, url, http.StatusFound)
				return
			}
		}

		data := templates.NewPageData2(100, 102, html) // index 102 for urlredirect
		renderer.RenderPage(w, "urlredirect/urlredirect.html", data)
	}
}
