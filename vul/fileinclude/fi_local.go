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
		filename := r.URL.Query().Get("file")
		if filename == "" {
			filename = "include/file1.php"
		}

		safeBase := "vul/fileinclude/include"
		cleanPath := filepath.Clean(filepath.Join(safeBase, filename))

		if !filepath.HasPrefix(cleanPath, safeBase) {
			data := templates.NewPageData2(74, 76, `<p style="color:red;">非法访问路径</p>`)
			renderer.RenderPage(w, "fileinclude/fi_local.html", data)
			return
		}

		content, err := os.ReadFile(cleanPath)
		html := ""
		if err != nil {
			html = fmt.Sprintf(`<p style="color:red;">无法读取文件: %s</p>`, err.Error())
		} else {
			html = fmt.Sprintf(`<h4>包含文件内容: %s</h4><pre>%s</pre>`, template.HTMLEscapeString(filename), template.HTMLEscapeString(string(content)))
		}

		data := templates.NewPageData2(74, 76, html)
		renderer.RenderPage(w, "fileinclude/fi_local.html", data)
	}
}
