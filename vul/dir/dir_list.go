package dir

import (
	"fmt"
	"html"
	"net/http"
	"os"
	"path/filepath"
	"pikachu-go/templates"
)

func DirListHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Query().Get("title")
		var htmlContent string

		// 默认路径拼接：soup/xxx.php
		soupPath := filepath.Join("templates/dir/soup", title+".html")

		if fileExists(soupPath) {
			// 如果是 soup 中合法模块
			data, err := os.ReadFile(soupPath)
			if err != nil {
				htmlContent = fmt.Sprintf(`<p style="color:red;">读取 soup 页面失败: %s</p>`, err.Error())
			} else {
				htmlContent = string(data)
			}
		} else {
			// 否则模拟任意文件包含漏洞
			// 支持文件路径如 ../../../../etc/passwd
			data, err := os.ReadFile(title)
			if err != nil {
				htmlContent = fmt.Sprintf(`<p style="color:red;">文件读取失败: %s</p>`, html.EscapeString(err.Error()))
			} else {
				htmlContent = "<pre>" + html.EscapeString(string(data)) + "</pre>"
			}
		}

		// 设置页面状态（左侧菜单高亮，面包屑导航）
		data := templates.NewPageData2(80, 82, htmlContent)
		err := renderer.RenderPage(w, "dir/dir_list.html", data)
		if err != nil {
			http.Error(w, "页面渲染失败："+err.Error(), http.StatusInternalServerError)
		}
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
