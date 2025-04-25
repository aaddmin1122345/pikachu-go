package burteforce

import (
	"html/template"
	"net/http"

	"pikachu-go/templates"
)

// BfClientHandler 客户端 JS 限制（服务端不校验）
func BfClientHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewPageData2(1, 5, "")
		if r.Method == "POST" {
			r.ParseForm()
			username := r.Form.Get("username")
			password := r.Form.Get("password")

			// 由于验证码验证在前端JS中完成，服务端不校验验证码
			if username == "" || password == "" {
				data.HtmlMsg = template.HTML("<p>please input username and password～</p>")
			} else if username == "admin" && password == "123456" {
				data.HtmlMsg = template.HTML("<p>login success</p>")
			} else {
				data.HtmlMsg = template.HTML("<p>username or password is not exists～</p>")
			}
		}

		renderer.RenderPage(w, "burteforce/bf_client.html", data)
	}
}
