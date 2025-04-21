package csrftoken

import (
	"net/http"
	"pikachu-go/templates"
)

// TokenGetHandler 显示用户信息与带 token 的链接
func TokenGetHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _ := r.Cookie("token_get_user")
		nick, _ := r.Cookie("token_get_nick")
		token, _ := r.Cookie("token_get_token")

		html := `
		<p>你好：` + user.Value + `（` + nick.Value + `）</p>
		<p><a href='/vul/csrf/csrftoken/token_get_edit?nickname=test&token=` + token.Value + `'>修改昵称（带 token）</a></p>
		`
		data := templates.NewPageData2(2, 16, html)
		renderer.RenderPage(w, "csrf/csrftoken/token_get.html", data)
	}
}
