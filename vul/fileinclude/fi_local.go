package fileinclude

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"pikachu-go/templates"
)

// FiLocalHandler 模拟本地文件包含漏洞
func FiLocalHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		html := ""

		if r.URL.Query().Get("submit") != "" && r.URL.Query().Get("filename") != "" {
			// 从URL获取文件名
			filename := r.URL.Query().Get("filename")

			// 直接包含文件，模拟PHP中的文件包含漏洞
			// 这里故意不做任何安全检查，以展示本地文件包含漏洞
			includePath := filepath.Join("vul/fileinclude/include", filename)
			content, err := os.ReadFile(includePath)

			if err != nil {
				html = fmt.Sprintf("<p style='color:red;'>无法读取文件: %s</p>", template.HTMLEscapeString(err.Error()))
			} else {
				html = string(content)
			}
		}

		data := templates.NewPageData2(55, 57, html)
		renderer.RenderPage(w, "fileinclude/fi_local.html", data)
	}
}
