package ssrf

import (
	"fmt"
	"io"
	"net/http"
	"pikachu-go/templates"
)

// SsrfCurlHandler 使用 net/http 模拟 curl SSRF
func SsrfCurlHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result := ""
		if r.Method == http.MethodPost {
			target := r.FormValue("url")
			if target != "" {
				resp, err := http.Get(target)
				if err != nil {
					result = fmt.Sprintf(`<p style="color:red;">请求失败：%s</p>`, err.Error())
				} else {
					defer resp.Body.Close()
					body, _ := io.ReadAll(resp.Body)
					result = fmt.Sprintf(`<h4>响应内容：</h4><pre>%s</pre>`, body)
				}
			}
		}

		data := templates.NewPageData2(105, 107, result)
		renderer.RenderPage(w, "ssrf/ssrf_curl.html", data)
	}
}
