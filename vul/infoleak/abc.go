package infoleak

import (
	"fmt"
	"net/http"
	"pikachu-go/templates"
)

// ABCHandler 处理abc页面的逻辑
func ABCHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 调试模式：直接查看页面，跳过验证
		if r.URL.Query().Get("debug") == "1" {
			fmt.Println("调试模式：直接显示页面")
			data := templates.NewPageData2(85, 87, "")
			renderer.RenderPage(w, "infoleak/abc.html", data)
			return
		}

		// 退出逻辑
		if r.URL.Query().Get("logout") == "1" {
			http.SetCookie(w, &http.Cookie{Name: "abc_uname", Value: "", Path: "/"})
			http.SetCookie(w, &http.Cookie{Name: "abc_pw", Value: "", Path: "/"})
			http.Redirect(w, r, "/vul/infoleak/findabc", http.StatusSeeOther)
			return
		}

		// 读取cookie验证身份
		var username, password string
		var isAuthenticated bool = false

		// 打印调试信息
		fmt.Println("当前请求的所有Cookie:")
		for _, cookie := range r.Cookies() {
			fmt.Printf("Cookie: %s = %s\n", cookie.Name, cookie.Value)
			if cookie.Name == "abc_uname" {
				username = cookie.Value
			}
			if cookie.Name == "abc_pw" {
				password = cookie.Value
			}
		}

		fmt.Printf("提取到的用户信息: username=%s, password=%s\n", username, password)

		// 简单验证，只要cookie存在即可访问
		if username == "lili" && password == "123456" {
			isAuthenticated = true
			fmt.Println("验证通过，允许访问")
		} else {
			fmt.Println("验证失败，拒绝访问")
		}

		// 若未通过验证，则使用标准HTML警告用户未授权
		if !isAuthenticated {
			htmlMsg := `<p class="notice">未授权访问！请通过正确途径登录。</p><p class="notice">提示：请使用lili/123456账号登录</p>`
			data := templates.NewPageData2(85, 87, htmlMsg)
			renderer.RenderPage(w, "infoleak/findabc.html", data)
			return
		}

		// 已授权，显示内容
		data := templates.NewPageData2(85, 87, "")
		renderer.RenderPage(w, "infoleak/abc.html", data)
	}
}
