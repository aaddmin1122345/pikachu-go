package dir

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"pikachu-go/templates"
	"strings"
)

// DirHandler 目录遍历（模拟目录遍历漏洞）
func DirHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqPath := r.URL.Query().Get("path")
		if reqPath == "" {
			reqPath = "." // 默认当前目录
		}
		// 简单过滤避免穿越
		cleanPath := filepath.Clean(reqPath)

		entries, err := os.ReadDir(cleanPath)
		if err != nil {
			data := templates.NewPageData2(4, 18, "读取目录失败："+template.HTMLEscapeString(err.Error()))
			renderer.RenderPage(w, "dir/dir.html", data)
			return
		}

		// 构造 HTML 表格
		var builder strings.Builder
		builder.WriteString(`<p>当前路径：<strong>` + template.HTMLEscapeString(cleanPath) + `</strong></p>`)
		builder.WriteString(`<table border="1" cellpadding="5"><tr><th>名称</th><th>类型</th></tr>`)

		for _, entry := range entries {
			name := entry.Name()
			entryType := "文件"
			if entry.IsDir() {
				entryType = "目录"
			}
			builder.WriteString(`<tr><td>` + template.HTMLEscapeString(name) + `</td><td>` + entryType + `</td></tr>`)
		}
		builder.WriteString(`</table>`)

		data := templates.NewPageData2(4, 18, builder.String())
		renderer.RenderPage(w, "dir/dir.html", data)
	}
}
