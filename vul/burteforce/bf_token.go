package burteforce

import (
	"net/http"
	"strconv"
	"time"

	"pikachu-go/templates"
)

func BfTokenHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, _ := sessionStore.Get(r, "pikachu-session")
		var msg string

		if r.Method == http.MethodPost {
			username := r.FormValue("username")
			password := r.FormValue("password")
			formToken := r.FormValue("token")
			sessToken, _ := sess.Values["csrf_token"].(string)

			switch {
			case username == "":
				msg = `<p class="notice">用户名不能为空</p>`
			case password == "":
				msg = `<p class="notice">密码不能为空</p>`
			case formToken == "":
				msg = `<p class="notice">令牌不能为空</p>`
			case formToken != sessToken:
				msg = `<p class="notice">令牌验证失败</p>`
			case username == "admin" && password == "123456":
				msg = `<p class="notice">登录成功！</p>`
			default:
				msg = `<p class="notice">用户名或密码错误</p>`
			}
		}

		// 生成新 token
		token := "token_" + strconv.FormatInt(time.Now().Unix()/100*100, 10)
		sess.Values["csrf_token"] = token
		_ = sess.Save(r, w)

		data := templates.NewPageData2(1, 6, msg)
		data.Extra["Token"] = token
		renderer.RenderPage(w, "burteforce/bf_token.html", data)
	}
}
