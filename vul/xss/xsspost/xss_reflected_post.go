package xsspost

import (
	"fmt"
	"html/template"
	"net/http"
	"pikachu-go/templates"
)

func XssReflectedPostHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewPageData2(7, 10, "")

		// 处理退出登录
		if r.URL.Query().Get("logout") == "1" {
			http.SetCookie(w, &http.Cookie{Name: "ant_uname", Value: "", MaxAge: -1})
			http.SetCookie(w, &http.Cookie{Name: "ant_pw", Value: "", MaxAge: -1})
			http.Redirect(w, r, "/vul/xss/xsspost/post_login", http.StatusFound)
			return
		}

		// 验证登录状态
		_, err := r.Cookie("ant_uname")
		if err != nil {
			http.Redirect(w, r, "/vul/xss/xsspost/post_login", http.StatusFound)
			return
		}

		// 设置登录状态信息
		data.Extra["state"] = template.HTML("你已经登陆成功,<a href=\"/vul/xss/xsspost/xss_reflected_post?logout=1\">退出登陆</a>")

		// 处理POST提交
		if r.Method == http.MethodPost {
			message := r.FormValue("message")
			if message == "" {
				data.HtmlMsg = template.HTML("<p class='notice'>输入'kobe'试试-_-</p>")
			} else if message == "kobe" {
				data.HtmlMsg = template.HTML(fmt.Sprintf("<p class='notice'>愿你和%s一样，永远年轻，永远热血沸腾！</p><img src='/assets/images/nbaplayer/kobe.png' />", message))
			} else {
				// XSS漏洞点：直接输出用户输入而不进行过滤
				data.HtmlMsg = template.HTML(fmt.Sprintf("<p class='notice'>who is %s,i don't care!</p>", message))
			}
		}

		renderer.RenderPage(w, "xss/xsspost/xss_reflected_post.html", data)
	}
}
