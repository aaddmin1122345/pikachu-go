package urlredirect

import (
	"net/http"
	"pikachu-go/templates"
)

func UnsafeReHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url != "" {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(`
				<html>
				<head><title>跳转中</title></head>
				<body>
					<p>跳转中，请稍后...</p>
					<script>
						window.location.href = "` + url + `";
					</script>
				</body>
				</html>
			`))
			return
		}

		data := templates.NewPageData2(100, 101, "")
		renderer.RenderPage(w, "urlredirect/unsafere.html", data)
	}
}
