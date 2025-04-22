package csrfpost

import (
	"net/http"
	"pikachu-go/templates"
)

// CsrfPostHandler 展示昵称 + 提交表单（可被伪造）
func CsrfPostHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _ := r.Cookie("csrf_post_user")
		nick, _ := r.Cookie("csrf_post_nick")

		html := `
		<p>你好：` + user.Value + `（` + nick.Value + `）</p>
		<form action="/vul/csrf/csrfpost/csrf_post_edit" method="post">
			新的昵称：<input type="text" name="nickname"><br>
			<input type="submit" value="提交修改">
		</form>
		`
		data := templates.NewPageData2(25, 28, html)
		renderer.RenderPage(w, "csrf/csrfpost/csrf_post.php", data)
	}
}
