package burteforce

import (
	"html/template"
	"net/http"

	"pikachu-go/templates"
)

// BfFormHandler 普通登录表单（无防护）
func BfFormHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewPageData2(1, 3, "")

		if r.Method == http.MethodPost {
			username := r.FormValue("username")
			password := r.FormValue("password")

			switch {
			case username == "":
				data.HtmlMsg = template.HTML("<p>please input username～</p>")
			case password == "":
				data.HtmlMsg = template.HTML("<p>please input password～</p>")
			case username == "admin" && password == "123456":
				data.HtmlMsg = template.HTML("<p>login success</p>")
			default:
				data.HtmlMsg = template.HTML("<p>username or password is not exists～</p>")
			}
		}

		renderer.RenderPage(w, "burteforce/bf_form.html", data)
	}
}
