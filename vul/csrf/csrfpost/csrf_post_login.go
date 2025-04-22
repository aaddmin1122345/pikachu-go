package csrfpost

import (
	"net/http"
	"pikachu-go/templates"
)

// CsrfPostLoginHandler 模拟登录并设置 cookie
func CsrfPostLoginHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := ""
		if r.Method == http.MethodPost {
			user := r.FormValue("username")
			nick := r.FormValue("nickname")
			http.SetCookie(w, &http.Cookie{Name: "csrf_post_user", Value: user, Path: "/"})
			http.SetCookie(w, &http.Cookie{Name: "csrf_post_nick", Value: nick, Path: "/"})
			msg = "登录成功"
		}
		data := templates.NewPageData2(25, 28, msg)
		renderer.RenderPage(w, "csrf/csrfpost/csrf_post_login.php", data)
	}
}
