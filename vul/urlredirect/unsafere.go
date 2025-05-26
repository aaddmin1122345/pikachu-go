package urlredirect

import (
	"html/template"
	"net/http"
	"pikachu-go/templates"
	"strings"
)

// UnsafeReHandler 处理基于JavaScript的不安全重定向
func UnsafeReHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url != "" {
			// 确保URL字符串安全，防止XSS
			if !strings.HasPrefix(url, "/") && !strings.Contains(url, "://") {
				url = "/" + url
			}

			// 使用JavaScript进行客户端重定向
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			html := `
				<html>
				<head>
					<title>跳转中</title>
					<meta charset="utf-8">
				</head>
				<body>
					<p>跳转中，请稍后...</p>
					<script>
						window.location.href = "` + template.JSEscapeString(url) + `";
					</script>
				</body>
				</html>
			`
			w.Write([]byte(html))
			return
		}

		data := templates.NewPageData2(100, 101, "")
		renderer.RenderPage(w, "urlredirect/unsafere.html", data)
	}
}
