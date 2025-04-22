package infoleak

import (
	"fmt"
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
			http.Redirect(w, r, "/vul/infoleak/findabc", http.StatusSeeOther)
			return
		}

		// 模拟信息泄露：读取 cookie
		var username, password string
		var hasRequiredCookies bool = false

		for _, cookie := range r.Cookies() {
			if cookie.Name == "abc[uname]" {
				username = cookie.Value
			}
			if cookie.Name == "abc[pw]" {
				password = cookie.Value
			}
		}

		// 检查是否具有正确的cookie访问权限
		if username == "admin" && password == "123456" {
			hasRequiredCookies = true
		}

		// 如果没有正确的cookie，返回错误页面
		if !hasRequiredCookies {
			htmlMsg := `
				<div class="pika-error" style="padding: 15px; border-radius: 4px;">
					<h4 style="margin-top: 0;"><i class="ace-icon fa fa-exclamation-circle"></i> 访问受限</h4>
					<p>您没有查看此页面的权限。请通过正确途径获取访问权限。</p>
					<p>提示：您可能需要先设置正确的cookie值...</p>
				</div>
			`
			data := templates.NewPageData2(85, 87, htmlMsg)
			renderer.RenderPage(w, "infoleak/abc.html", data)
			return
		}

		// 如果具有正确的cookie值，显示成功页面
		var leakInfo []string
		for _, cookie := range r.Cookies() {
			if strings.HasPrefix(cookie.Name, "abc[") {
				leakInfo = append(leakInfo, cookie.Name+" = "+cookie.Value)
			}
		}

		cookieInfo := ""
		if len(leakInfo) > 0 {
			cookieInfo = strings.Join(leakInfo, "\n")
		}

		htmlMsg := fmt.Sprintf(`
			<div style="display:none">
				<!-- Cookie信息：%s -->
				<!-- 这是一个隐藏的信息泄露演示 -->
			</div>
		`, cookieInfo)

		data := templates.NewPageData2(85, 87, htmlMsg)
		renderer.RenderPage(w, "infoleak/abc.html", data)
	}
}
