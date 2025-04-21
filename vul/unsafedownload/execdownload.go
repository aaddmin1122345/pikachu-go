package unsafedownload

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// ExecDownloadHandler 处理文件下载请求，模拟漏洞：无验证 filename 参数
func ExecDownloadHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := r.URL.Query().Get("filename")
		if filename == "" || strings.Contains(filename, "..") {
			http.Error(w, "非法文件名", http.StatusBadRequest)
			return
		}

		filePath := filepath.Join("vul/unsafedownload/download", filename)

		// 检查文件是否存在
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			http.Error(w, "文件不存在", http.StatusNotFound)
			return
		}

		// 设置头部
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")

		// 打开文件并传输
		f, err := os.Open(filePath)
		if err != nil {
			http.Error(w, "文件打开失败", http.StatusInternalServerError)
			return
		}
		defer f.Close()

		io.Copy(w, f)
	}
}
