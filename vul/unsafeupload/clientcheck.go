package unsafeupload

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

// ClientcheckHandler 前端JS校验的文件上传漏洞实现
func ClientcheckHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		htmlMsg := ""

		if r.Method == http.MethodPost {
			// 处理文件上传
			file, header, err := r.FormFile("uploadfile")
			if err != nil {
				htmlMsg = "<p class='notice'>上传失败，请选择文件</p>"
			} else {
				defer file.Close()

				// 确保上传目录存在
				uploadDir := "vul/unsafeupload/uploads"
				if err := os.MkdirAll(uploadDir, 0755); err != nil {
					htmlMsg = fmt.Sprintf("<p class='notice'>创建目录失败: %s</p>", err.Error())
				} else {
					// 使用原始文件名保存，这是不安全的做法，容易导致安全问题
					// 该功能故意不做服务器端校验，只依赖客户端JS校验
					filename := header.Filename
					filepath := filepath.Join(uploadDir, filename)

					// 创建目标文件
					dst, err := os.Create(filepath)
					if err != nil {
						htmlMsg = fmt.Sprintf("<p class='notice'>创建文件失败: %s</p>", err.Error())
					} else {
						defer dst.Close()

						// 将上传的文件内容写入目标文件
						if _, err := io.Copy(dst, file); err != nil {
							htmlMsg = fmt.Sprintf("<p class='notice'>保存文件失败: %s</p>", err.Error())
						} else {
							// 文件上传成功，返回成功信息和文件路径
							filePath := strings.ReplaceAll(filepath, "\\", "/") // 统一路径分隔符
							htmlMsg = fmt.Sprintf("<p class='notice'>文件上传成功!</p>"+
								"<p class='notice'>文件保存的路径为：%s</p>"+
								"<p class='notice'>如果是图片，<a href='/%s' target='_blank'>点击浏览</a></p>",
								filePath, filePath)

							// 添加提示信息
							if !strings.Contains(strings.ToLower(filename), ".php") {
								htmlMsg += "<p class='notice' style='color:blue;'>提示: 尝试上传PHP文件绕过客户端校验</p>"
							}
						}
					}
				}
			}
		}

		data := templates.PageData{
			Active:  make([]string, 130),
			HtmlMsg: template.HTML(htmlMsg),
		}
		data.Active[65] = "active open"
		data.Active[67] = "active"

		renderer.RenderPage(w, "unsafeupload/clientcheck.html", data)
	}
}
