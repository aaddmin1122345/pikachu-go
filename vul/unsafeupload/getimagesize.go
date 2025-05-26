package unsafeupload

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"  // 注册GIF解码器
	_ "image/jpeg" // 注册JPEG解码器
	_ "image/png"  // 注册PNG解码器
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"pikachu-go/templates"
	"strings"
	"time"
)

// GetimagesizeHandler 实现图像尺寸和有效性验证的上传处理
func GetimagesizeHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := ""

		if r.Method == http.MethodPost {
			file, header, err := r.FormFile("uploadfile")
			if err == nil {
				defer file.Close()

				// 读取文件内容到内存中
				fileBytes, err := ioutil.ReadAll(file)
				if err != nil {
					msg = "<p class='notice'>读取文件失败</p>"
				} else {
					// 验证图片的有效性
					_, format, err := image.DecodeConfig(bytes.NewReader(fileBytes))
					if err != nil {
						msg = "<p class='notice'>上传失败: 文件不是有效的图片</p>"
					} else {
						// 检查MIME类型
						contentType := header.Header.Get("Content-Type")
						allowedTypes := map[string]string{
							"image/jpeg": "jpeg",
							"image/png":  "png",
							"image/gif":  "gif",
						}

						// 校验格式与内容类型是否匹配
						expectedFormat, isAllowedType := allowedTypes[contentType]
						if !isAllowedType {
							msg = "<p class='notice'>上传失败: 不允许的文件类型</p>"
						} else if expectedFormat != format && !(expectedFormat == "jpeg" && format == "jpg") {
							msg = "<p class='notice'>上传失败: 文件格式与内容类型不匹配</p>"
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

								// 写入文件
								err = ioutil.WriteFile(dst, fileBytes, 0644)
								if err != nil {
									msg = "<p class='notice'>保存文件失败</p>"
								} else {
									filepath := strings.ReplaceAll(dst, "\\", "/") // 统一路径分隔符
									msg = fmt.Sprintf("<p class='notice'>文件上传成功</p>"+
										"<p class='notice'>文件保存的路径为：%s</p>"+
										"<p class='notice'>文件格式: %s</p>"+
										"<p class='notice'><a href='/%s' target='_blank'>点击浏览</a></p>",
										filepath, format, filepath)
								}
							}
						}
					}
				}
			} else {
				msg = "<p class='notice'>上传失败, 请选择文件</p>"
			}
		}

		data := templates.NewPageData2(65, 69, msg)
		renderer.RenderPage(w, "unsafeupload/getimagesize.html", data)
	}
}
