package csrfget

import (
	"net/http"
	"pikachu-go/templates"
)

// CsrfGetHandler 显示用户信息页面（从 cookie 中读取）
func CsrfGetHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := ""
		nick := ""
		if u, err := r.Cookie("csrf_get_user"); err == nil {
			user = u.Value
		}
		if n, err := r.Cookie("csrf_get_nick"); err == nil {
			nick = n.Value
		}

		html := `
		<p>你好：` + user + `(` + nick + `)</p>
		<p><a href='/vul/csrf/csrfget/csrf_get_edit?nickname=test'>修改昵称</a></p>
		`
		data := templates.NewPageData2(2, 14, html)
		renderer.RenderPage(w, "csrf/csrfget/csrf_get.php", data)
	}
}
