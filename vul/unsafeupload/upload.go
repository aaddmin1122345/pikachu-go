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

// UploadHandler 渲染上传页面 + 处理上传逻辑
// 这是漏洞概述页面的处理器，保留一些基本的安全措施
func UploadHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := ""

		if r.Method == http.MethodPost {
			file, header, err := r.FormFile("upload_file")
			if err != nil {
				msg = `<p style="color:red;">上传失败：未选择文件</p>`
			} else {
				defer file.Close()

				// 确保上传目录存在
				uploadDir := "vul/unsafeupload/uploads"
				if err := os.MkdirAll(uploadDir, 0755); err != nil {
					msg = `<p style="color:red;">上传失败：无法创建上传目录</p>`
				} else {
					// 添加时间戳前缀，避免文件名冲突
					timestamp := time.Now().Unix()
					safeFilename := fmt.Sprintf("%d_%s", timestamp, header.Filename)
					savePath := filepath.Join(uploadDir, safeFilename)

					out, err := os.Create(savePath)
					if err != nil {
						msg = `<p style="color:red;">上传失败：无法保存文件</p>`
					} else {
						defer out.Close()
						_, err = io.Copy(out, file)
						if err != nil {
							msg = `<p style="color:red;">上传失败：保存文件时发生错误</p>`
						} else {
							// 统一路径分隔符
							filePath := strings.ReplaceAll(savePath, "\\", "/")
							msg = fmt.Sprintf(`<p style="color:green;">上传成功！你可以访问：<a href="/%s" target="_blank">%s</a></p>`,
								filePath, header.Filename)

							// // 添加关于其他示例的说明
							// msg += `<div style="margin-top:20px;">
							// 	<h4>文件上传漏洞示例:</h4>
							// 	<ul>
							// 		<li><a href="/vul/unsafeupload/clientcheck">客户端JS校验</a> - 演示前端验证的缺陷</li>
							// 		<li><a href="/vul/unsafeupload/servercheck">服务端MIME类型校验</a> - 演示基于MIME类型的验证</li>
							// 		<li><a href="/vul/unsafeupload/getimagesize">图像有效性校验</a> - 演示更严格的图像验证</li>
							// 	</ul>
							// </div>`
						}
					}
				}
			}
		}

		data := templates.NewPageData2(65, 66, msg)
		renderer.RenderPage(w, "unsafeupload/upload.html", data)
	}
}
