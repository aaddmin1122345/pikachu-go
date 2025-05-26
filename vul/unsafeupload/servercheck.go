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

// ServercheckHandler 实现服务端MIME类型检查的上传处理
func ServercheckHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := ""

		if r.Method == http.MethodPost {
			file, header, err := r.FormFile("uploadfile")
			if err == nil {
				defer file.Close()

				// 检查MIME类型
				contentType := header.Header.Get("Content-Type")
				allowedTypes := map[string]string{
					"image/jpeg": ".jpg",
					"image/png":  ".png",
					"image/gif":  ".gif",
				}

				fileExt := strings.ToLower(filepath.Ext(header.Filename))

				// 检查是否是允许的MIME类型
				expectedExt, isAllowedType := allowedTypes[contentType]

				if !isAllowedType {
					msg = "<p class='notice'>上传失败: 不允许的文件类型</p>"
				} else if expectedExt != fileExt && ".jpeg" != fileExt { // 特殊处理.jpeg扩展名
					// 检查文件扩展名是否与MIME类型匹配
					msg = "<p class='notice'>上传失败: 文件扩展名与内容类型不匹配</p>"
				} else {
					// 确保目录存在
					saveDir := "vul/unsafeupload/uploads"
					err := os.MkdirAll(saveDir, 0755)
					if err != nil {
						msg = "<p class='notice'>创建目录失败</p>"
					} else {
						// 为文件名添加时间戳前缀以避免文件覆盖
						timestamp := time.Now().Unix()
						filename := fmt.Sprintf("%d_%s", timestamp, header.Filename)
						dst := filepath.Join(saveDir, filename)

						out, err := os.Create(dst)
						if err == nil {
							defer out.Close()
							_, err := io.Copy(out, file)
							if err != nil {
								msg = "<p class='notice'>保存文件失败</p>"
							} else {
								filepath := strings.ReplaceAll(dst, "\\", "/") // 统一路径分隔符
								msg = fmt.Sprintf("<p class='notice'>文件上传成功</p>"+
									"<p class='notice'>文件保存的路径为：%s</p>"+
									"<p class='notice'>如果是图片，<a href='/%s' target='_blank'>点击浏览</a></p>",
									filepath, filepath)
							}
						} else {
							msg = "<p class='notice'>创建文件失败</p>"
						}
					}
				}
			} else {
				msg = "<p class='notice'>上传失败, 请选择文件</p>"
			}
		}

		data := templates.NewPageData2(65, 68, msg)
		renderer.RenderPage(w, "unsafeupload/servercheck.html", data)
	}
}
