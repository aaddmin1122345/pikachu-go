package xss

import (
	"fmt"
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
					// 用于href属性中，替换尖括号
					if message == "" {
						data.HtmlMsg = template.HTML("<p class='notice'>叫你输入个url,你咋不听?</p>")
					} else if message == "www.baidu.com" {
						data.HtmlMsg = template.HTML("<p class='notice'>我靠,我真想不到你是这样的一个人</p>")
					} else {
						// 这里应该使用htmlspecialchars编码后再输出到href属性
						message = strings.ReplaceAll(strings.ReplaceAll(message, "<", "["), ">", "]")
						data.HtmlMsg = template.HTML(fmt.Sprintf("<a href='%s'> 阁下自己输入的url还请自己点一下吧</a>", message))
					}
				case "xss_04":
					// JS输出场景，直接将用户输入嵌入JS中
					// 这里应将用户输入进行JavaScript转义
					jsvar := message
					imgTag := ""

					if message == "tmac" {
						imgTag = "<img src='/assets/images/nbaplayer/tmac.jpeg' />"
					}

					// 添加JS变量到Extra
					data.Extra = make(map[string]interface{})
					data.Extra["Jsvar"] = jsvar
					data.HtmlMsg = template.HTML(imgTag)
				}
			}
		}

		renderer.RenderPage(w, "xss/"+variant+".html", data)
	}
}
