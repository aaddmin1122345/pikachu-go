package ssrf

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"pikachu-go/templates"
	"strings"
)

// SsrfFgcHandler 使用 file_get_contents 模拟 SSRF
func SsrfFgcHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 优先处理GET参数，模拟PHP版本行为
		fileParam := r.URL.Query().Get("file")

		if fileParam != "" {
			// 直接处理文件请求并显示结果
			content, err := fetchFileContent(fileParam)
			if err != nil {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				result := fmt.Sprintf(`<p style="color:red;">请求失败：%s</p>`, err.Error())
				w.Write([]byte(result))
			} else {
				// 直接输出内容而不是渲染模板
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.Write([]byte(content))
			}
			return
		}

		// 构建示例链接 - 完全匹配PHP版本的格式
		exampleLink := `<a href="/vul/ssrf/ssrf_fgc?file=http://127.0.0.1:8888/vul/ssrf/ssrf_info/info2">反正都读了,那就在来一首吧</a>`

		data := templates.NewPageData2(105, 108, exampleLink)
		renderer.RenderPage(w, "ssrf/ssrf_fgc.html", data)
	}
}

// fetchFileContent 模拟PHP的file_get_contents函数
func fetchFileContent(path string) (string, error) {
	// 支持PHP过滤器协议(模拟)
	if strings.HasPrefix(path, "php://") {
		// 这里只是模拟，实际无法完全支持php://filter功能
		return "PHP Filter Protocol is not fully supported in Go version.", nil
	}

	// 处理文件协议
	if strings.HasPrefix(path, "file://") {
		localPath := strings.TrimPrefix(path, "file://")
		content, err := os.ReadFile(localPath)
		if err != nil {
			return "", err
		}
		return string(content), nil
	}

	// 处理gopher协议
	if strings.HasPrefix(path, "gopher://") {
		return "Gopher Protocol is not fully supported in Go version.", nil
	}

	// 处理dict协议
	if strings.HasPrefix(path, "dict://") {
		return "Dict Protocol is not fully supported in Go version.", nil
	}

	// 处理HTTP/HTTPS请求
	resp, err := http.Get(path)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
