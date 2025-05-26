package unsafeupload

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"pikachu-go/templates"
	"strings"
	"time"
)

// ClientcheckHandler 前端JS校验的文件上传漏洞实现
// 该功能有意保留一些漏洞，但增加基本的安全措施
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
					// 添加基本的文件扩展名检查
					filename := header.Filename
					lowerFilename := strings.ToLower(filename)

					// 虽然前端有JS校验，但我们这里也做一个基本的后端验证
					// 这个验证故意留有漏洞，允许绕过（例如双扩展名攻击）
					forbiddenExtensions := []string{".php", ".asp", ".aspx", ".exe", ".sh", ".bat"}
					isForbidden := false

					for _, ext := range forbiddenExtensions {
						if strings.HasSuffix(lowerFilename, ext) {
							isForbidden = true
							break
						}
					}

					if isForbidden {
						htmlMsg = "<p class='notice'>上传失败: 不允许上传此类型的文件</p>"
					} else {
						// 为文件名添加时间戳前缀以避免文件覆盖
						timestamp := time.Now().Unix()
						safeFilename := fmt.Sprintf("%d_%s", timestamp, filename)
						filepath := filepath.Join(uploadDir, safeFilename)

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

								// // 添加提示信息
								// if !strings.Contains(strings.ToLower(filename), ".php") {
								// 	htmlMsg += "<p class='notice' style='color:blue;'>提示: 尝试上传PHP文件绕过客户端校验</p>"
								// }
							}
						}
					}
				}
			}
		}

		data := templates.NewPageData2(65, 67, htmlMsg)
		renderer.RenderPage(w, "unsafeupload/clientcheck.html", data)
	}
}
