package burteforce

import (
	"html/template"
	"net/http"
	"strings"

	"pikachu-go/templates"
)

// BfServerHandler 服务端控制爆破尝试 + 验证码校验
func BfServerHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, _ := sessionStore.Get(r, "pikachu-session")
		data := templates.NewPageData2(1, 4, "")

		// 漏洞点：不主动刷新验证码，只有在用户点击验证码图片时才会刷新
		// 这样可以用于模拟绕过验证码的暴力破解

		if r.Method == http.MethodPost {
			username := r.FormValue("username")
			password := r.FormValue("password")
			formVcode := r.FormValue("vcode")
			sessVcode, _ := sess.Values["vcode"].(string)

			switch {
			case username == "":
				data.HtmlMsg = template.HTML("<p>please input username～</p>")
			case password == "":
				data.HtmlMsg = template.HTML("<p>please input password～</p>")
			case formVcode == "":
				data.HtmlMsg = template.HTML("<p>please input verification code～</p>")
			case strings.ToLower(formVcode) != strings.ToLower(sessVcode):
				data.HtmlMsg = template.HTML("<p>verification code error～</p>")
			case username == "admin" && password == "123456":
				data.HtmlMsg = template.HTML("<p>login success</p>")
			default:
				data.HtmlMsg = template.HTML("<p>username or password is not exists～</p>")
			}
		}

		data.Extra["VcodeURL"] = "/vul/burteforce/vcode"
		renderer.RenderPage(w, "burteforce/bf_server.html", data)
	}
}
