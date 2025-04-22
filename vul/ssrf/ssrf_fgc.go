package ssrf

import (
	"fmt"
	"io"
	"net/http"
	"pikachu-go/templates"
)

// SsrfFgcHandler 使用 file_get_contents 模拟 SSRF（也使用 net/http）
func SsrfFgcHandler(renderer templates.Renderer) http.HandlerFunc {
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

		data := templates.NewPageData2(105, 108, result)
		renderer.RenderPage(w, "ssrf/ssrf_fgc.html", data)
	}
}
