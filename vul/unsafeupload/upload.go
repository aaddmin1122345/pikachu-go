package unsafeupload

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"pikachu-go/templates"
)

// UploadHandler 渲染上传页面 + 处理上传逻辑
func UploadHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := ""

		if r.Method == http.MethodPost {
			file, header, err := r.FormFile("upload_file")
			if err != nil {
				msg = `<p style="color:red;">上传失败：未选择文件</p>`
			} else {
				defer file.Close()

				savePath := filepath.Join("vul/unsafeupload/upload", header.Filename)
				out, err := os.Create(savePath)
				if err != nil {
					msg = `<p style="color:red;">上传失败：无法保存文件</p>`
				} else {
					defer out.Close()
					_, _ = io.Copy(out, file)
					msg = fmt.Sprintf(`<p style="color:green;">上传成功！你可以访问：<a href="/vul/unsafeupload/upload/%s" target="_blank">%s</a></p>`, header.Filename, header.Filename)
				}
			}
		}

		data := templates.NewPageData2(65, 66, msg)
		renderer.RenderPage(w, "unsafeupload/upload.html", data)
	}
}
