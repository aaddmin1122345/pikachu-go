package xss

import (
	"html/template"
	"net/http"
	"pikachu-go/templates"
	"regexp"
	"strings"
)

// RenderXssVariant 处理不同的XSS变种请求
func RenderXssVariant(renderer templates.Renderer, variant string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewPageData2(7, 0, "")

		// 设置不同变种对应的导航ID
		switch variant {
		case "xss_01":
			data = templates.NewPageData2(7, 15, "")
		case "xss_02":
			data = templates.NewPageData2(7, 16, "")
		case "xss_03":
			data = templates.NewPageData2(7, 17, "")
		case "xss_04":
			data = templates.NewPageData2(7, 18, "")
		}

		// 处理提交的内容
		if r.URL.Query().Get("submit") != "" {
			message := r.URL.Query().Get("message")
			if message != "" {
				switch variant {
				case "xss_01":
					// 使用正则过滤掉<script标签
					re := regexp.MustCompile(`<(.*)s(.*)c(.*)r(.*)i(.*)p(.*)t`)
					message = re.ReplaceAllString(message, "")
					if message == "yes" {
						data.HtmlMsg = template.HTML("<p>那就去人民广场一个人坐一会儿吧!</p>")
					} else {
						data.HtmlMsg = template.HTML("<p>别说这些'" + message + "'的话,不要怕,就是干!</p>")
					}
				case "xss_02":
					// 使用htmlspecialchars过滤
					message = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(message, "&", "&amp;"), "<", "&lt;"), ">", "&gt;")
					data.HtmlMsg = template.HTML("<p>输入的内容是：" + message + "</p>")
				case "xss_03":
					// 替换尖括号
					message = strings.ReplaceAll(strings.ReplaceAll(message, "<", "["), ">", "]")
					data.HtmlMsg = template.HTML("<p>输入的内容是：" + message + "</p>")
				case "xss_04":
					// 使用htmlentities编码
					message = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(message, "&", "&amp;"), "<", "&lt;"), ">", "&gt;"), "\"", "&quot;")
					data.HtmlMsg = template.HTML("<p>输入的内容是：" + message + "</p>")
				}
			}
		}

		renderer.RenderPage(w, "xss/"+variant+".html", data)
	}
}
