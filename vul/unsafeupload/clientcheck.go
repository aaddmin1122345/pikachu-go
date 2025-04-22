package unsafeupload

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"pikachu-go/templates"
)

// ClientcheckHandler 对应 clientcheck.php
func ClientcheckHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := ""

		if r.Method == http.MethodPost {
			file, header, err := r.FormFile("uploadfile")
			if err == nil {
				defer file.Close()
				saveDir := "vul/unsafeupload/uploads"
				os.MkdirAll(saveDir, 0755)
				dst := filepath.Join(saveDir, header.Filename)
				out, err := os.Create(dst)
				if err == nil {
					defer out.Close()
					io.Copy(out, file)
					msg = fmt.Sprintf("<p class='notice'>文件上传成功</p><p class='notice'>文件保存的路径为：%s</p>", dst)
				} else {
					msg = "<p class='notice'>保存失败</p>"
				}
			} else {
				msg = "<p class='notice'>上传失败</p>"
			}
		}

		data := templates.NewPageData2(65, 67, msg)
		renderer.RenderPage(w, "unsafeupload/clientcheck.html", data)
	}
}
