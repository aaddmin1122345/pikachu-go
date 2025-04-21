package fileinclude

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"pikachu-go/templates"
)

// FiRemoteHandler 模拟远程文件包含（RFI）行为
func FiRemoteHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("file")
		html := ""

		if url != "" {
			resp, err := http.Get(url)
			if err != nil {
				html = fmt.Sprintf(`<p style="color:red;">远程请求失败：%s</p>`, err.Error())
			} else {
				defer resp.Body.Close()
				body, _ := io.ReadAll(resp.Body)
				html = fmt.Sprintf(`<h4>包含远程内容: %s</h4><pre>%s</pre>`,
					template.HTMLEscapeString(url),
					template.HTMLEscapeString(string(body)))
			}
		}

		data := templates.NewPageData2(74, 77, html)
		renderer.RenderPage(w, "fileinclude/fi_remote.html", data)
	}
}
