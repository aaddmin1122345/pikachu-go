package infoleak

import (
	"fmt"
	"net/http"
	"pikachu-go/templates"
)

// FindABCHandler 实现信息泄露模块中的findabc功能
func FindABCHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := ""

		// 获取GET参数
		username := r.URL.Query().Get("username")
		password := r.URL.Query().Get("password")
		submit := r.URL.Query().Get("submit")

		// 如果提交了表单
		if submit != "" && username != "" && password != "" {
			// 在实际应用中应该进行数据库验证，这里简化处理
			// 模拟用户lili的登录验证
			if username == "lili" && password == "123456" {
				// 设置cookie - 使用HTTP Cookie而不是特殊格式
				http.SetCookie(w, &http.Cookie{
					Name:  "abc_uname",
					Value: username,
					Path:  "/",
				})
				http.SetCookie(w, &http.Cookie{
					Name:  "abc_pw",
					Value: password, // 实际应用中应该存储哈希值
					Path:  "/",
				})

				// 添加调试信息
				fmt.Printf("设置cookie成功: abc_uname=%s, abc_pw=%s\n", username, password)

				// 重定向到abc页面
				http.Redirect(w, r, "/vul/infoleak/abc", http.StatusFound)
				return
			} else if username == "lili" {
				msg = "<p class='notice'>您输入的密码错误</p>"
			} else {
				msg = "<p class='notice'>您输入的账号错误</p>"
			}
		}

		// 渲染页面
		data := templates.NewPageData2(85, 87, msg)
		renderer.RenderPage(w, "infoleak/findabc.html", data)
	}
}
