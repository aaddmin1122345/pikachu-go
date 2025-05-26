package ssrf_info

import (
	"net/http"
	"pikachu-go/templates"
)

// Info1Handler 处理诗歌请求
func Info1Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		renderer, err := templates.NewTemplateRenderer()
		if err != nil {
			http.Error(w, "模板加载失败", http.StatusInternalServerError)
			return
		}

		// 创建空的PageData对象
		data := templates.NewPageData2(105, 107, "")

		// 使用RenderPage方法渲染完整页面（包括头部和尾部）
		err = renderer.RenderPage(w, "ssrf/ssrf_info/info1.html", data)
		if err != nil {
			http.Error(w, "渲染模板失败", http.StatusInternalServerError)
			return
		}
	}
}
