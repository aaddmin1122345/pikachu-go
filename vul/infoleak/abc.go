package infoleak

import (
	"net/http"
	"pikachu-go/templates"
	"strings"
)

// ABCHandler 模拟信息泄露页面
func ABCHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 退出逻辑
		if r.URL.Query().Get("logout") == "1" {
			http.SetCookie(w, &http.Cookie{Name: "abc[uname]", Value: "", Path: "/"})
			http.SetCookie(w, &http.Cookie{Name: "abc[pw]", Value: "", Path: "/"})
			http.Redirect(w, r, "/vul/infoleak/abc", http.StatusSeeOther)
			return
		}

		// 模拟信息泄露：读取 cookie
		var leaks []string
		for _, cookie := range r.Cookies() {
			if strings.HasPrefix(cookie.Name, "abc[") {
				leaks = append(leaks, cookie.Name+" = "+cookie.Value)
			}
		}

		htmlMsg := "Cookie 泄露信息：\n"
		if len(leaks) > 0 {
			htmlMsg += strings.Join(leaks, "\n")
		} else {
			htmlMsg += "（无）"
		}

		data := templates.NewPageData2(66, 68, htmlMsg)
		renderer.RenderPage(w, "infoleak/abc.html", data)
	}
}
