package ssrf

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"pikachu-go/templates"
	"strings"
)

// SsrfCurlHandler 使用 net/http 模拟 curl SSRF
func SsrfCurlHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlParam := r.URL.Query().Get("url")

		if urlParam != "" {
			// 直接处理URL请求并显示结果
			content, err := fetchURLContent(urlParam)
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
		exampleLink := `<a href="/vul/ssrf/ssrf_curl?url=http://127.0.0.1:8888/vul/ssrf/ssrf_info/info1">累了吧,来读一首诗吧</a>`

		data := templates.NewPageData2(105, 107, exampleLink)
		renderer.RenderPage(w, "ssrf/ssrf_curl.html", data)
	}
}

// fetchURLContent 获取URL内容，支持多种协议
func fetchURLContent(url string) (string, error) {
	// 处理file协议
	if strings.HasPrefix(url, "file://") {
		path := strings.TrimPrefix(url, "file://")
		content, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		return string(content), nil
	}

	// 处理gopher协议
	if strings.HasPrefix(url, "gopher://") {
		return "Gopher Protocol is not fully supported in Go version.", nil
	}

	// 处理dict协议
	if strings.HasPrefix(url, "dict://") {
		return "Dict Protocol is not fully supported in Go version.", nil
	}

	// 处理PHP过滤器协议(模拟)
	if strings.HasPrefix(url, "php://") {
		return "PHP Filter Protocol is not fully supported in Go version.", nil
	}

	// 处理HTTP/HTTPS协议
	resp, err := http.Get(url)
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
