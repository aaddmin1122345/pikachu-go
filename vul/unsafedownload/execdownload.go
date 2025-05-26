package unsafedownload

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// ExecDownloadHandler 处理文件下载请求，模拟漏洞：无验证 filename 参数
func ExecDownloadHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := r.URL.Query().Get("filename")

		basePath := "vul/unsafedownload/download"
		filePath := filepath.Join(basePath, filename)

		// 漏洞点：未校验是否跳出 basePath
		// 如访问：?filename=../../../../etc/passwd

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			http.Error(w, "文件不存在", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")

		f, err := os.Open(filePath)
		if err != nil {
			http.Error(w, "文件打开失败", http.StatusInternalServerError)
			return
		}
		defer f.Close()

		io.Copy(w, f)
	}
}
