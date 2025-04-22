package burteforce

import (
	"net/http"
	"pikachu-go/templates"
)

// BfFormHandler 普通登录表单（无防护）
func BfFormHandler(renderer templates.Renderer) http.HandlerFunc {
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

		data := templates.NewPageData2(1, 3, msg)
		renderer.RenderPage(w, "burteforce/bf_form.html", data)
	}
}
