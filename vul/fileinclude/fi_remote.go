package fileinclude

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"pikachu-go/templates"
	"strings"
)

// FiRemoteHandler 模拟远程文件包含（RFI）行为
func FiRemoteHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		html := ""

		// // 模拟PHP的设置检查提示
		// html1 := "<p style='color: red'>提示: 在PHP中需要开启allow_url_include才能测试远程文件包含漏洞</p>"
		// html2 := "<p style='color: red'>提示: 在PHP中需要开启allow_url_fopen才能测试远程文件包含漏洞</p>"

		// 检查是否提交并且文件名不为空
		if r.URL.Query().Get("submit") != "" && r.URL.Query().Get("filename") != "" {
			// 获取文件路径
			filePath := r.URL.Query().Get("filename")

			// 如果是本地文件路径（不包含 http:// 或 https://）
			if !strings.HasPrefix(filePath, "http://") && !strings.HasPrefix(filePath, "https://") {
				// 从路径中提取文件名
				fileName := filepath.Base(filePath)
				// 构建完整的本地文件路径
				fullPath := filepath.Join("vul/fileinclude/include", fileName)
				content, err := os.ReadFile(fullPath)
				if err != nil {
					html = fmt.Sprintf("<p style='color:red;'>无法读取文件: %s</p>", template.HTMLEscapeString(err.Error()))
				} else {
					html = string(content)
				}
			} else {
				// 处理远程URL
				resp, err := http.Get(filePath)
				if err != nil {
					html = fmt.Sprintf("<p style='color:red;'>远程文件包含失败: %s</p>", template.HTMLEscapeString(err.Error()))
				} else {
					defer resp.Body.Close()
					body, _ := io.ReadAll(resp.Body)
					html = string(body)
				}
			}
		}

		// 组合所有HTML内容
		finalHTML := fmt.Sprintf("%s", html)

		data := templates.NewPageData2(55, 58, finalHTML)
		renderer.RenderPage(w, "fileinclude/fi_remote.html", data)
	}
}
