package csrftoken

import (
	"math/rand"
	"net/http"
	"pikachu-go/templates"
	"strconv"
	"time"
)

// TokenGetLoginHandler 登录并生成 token
func TokenGetLoginHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := ""
		if r.Method == http.MethodPost {
			user := r.FormValue("username")
			nick := r.FormValue("nickname")

			token := generateToken()
			http.SetCookie(w, &http.Cookie{Name: "token_get_user", Value: user, Path: "/"})
			http.SetCookie(w, &http.Cookie{Name: "token_get_nick", Value: nick, Path: "/"})
			http.SetCookie(w, &http.Cookie{Name: "token_get_token", Value: token, Path: "/"})
			msg = "登录成功，已生成 token"
		}

		data := templates.NewPageData2(25, 29, msg)
		renderer.RenderPage(w, "csrf/csrftoken/token_get_login.php", data)
	}
	// token 可加盐处理，这里简化为随机数字
}

func generateToken() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.FormatInt(rand.Int63(), 16)
}
