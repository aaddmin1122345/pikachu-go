package csrfget

import (
	"net/http"
	"pikachu-go/templates"
)

// CsrfGetLoginHandler 设置用户名到 cookie 中
func CsrfGetLoginHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := ""
		if r.Method == http.MethodPost {
			user := r.FormValue("username")
			nick := r.FormValue("nickname")
			http.SetCookie(w, &http.Cookie{Name: "csrf_get_user", Value: user, Path: "/"})
			http.SetCookie(w, &http.Cookie{Name: "csrf_get_nick", Value: nick, Path: "/"})
			msg = "登录成功"
		}

		data := templates.NewPageData2(25, 27, msg)
		renderer.RenderPage(w, "csrf/csrfget/csrf_get_login.php", data)
	}
}
