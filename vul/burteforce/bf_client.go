package burteforce

import (
	"net/http"
	"pikachu-go/templates"
)

// BfClientHandler 使用前端 JS 限制次数（服务端不判断）
func BfClientHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := ""
		if r.Method == http.MethodPost {
			username := r.FormValue("username")
			password := r.FormValue("password")

			if username == "admin" && password == "123456" {
				msg = `<p style="color:green;">登录成功！</p>`
			} else {
				msg = `<p style="color:red;">用户名或密码错误</p>`
			}
		}

		data := templates.NewPageData2(12, 15, msg)
		renderer.RenderPage(w, "burteforce/bf_client.html", data)
	}
}
