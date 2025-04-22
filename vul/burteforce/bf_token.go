package burteforce

import (
	"math/rand"
	"net/http"
	"pikachu-go/templates"
	"pikachu-go/utils"
	"strconv"
	"strings"
	"time"
)

// BfTokenHandler 加入 Token 验证，模拟防止 CSRF 爆破
func BfTokenHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := ""
		sessionToken := getOrSetToken(w, r)

		if r.Method == http.MethodPost {
			formToken := r.FormValue("token")
			username := r.FormValue("username")
			password := r.FormValue("password")
			vcode := r.FormValue("vcode")

			if formToken != sessionToken {
				msg = `<p style="color:red;">非法请求，Token 无效！</p>`
			} else if vcode == "" {
				msg = `<p style="color:red;">验证码不能为空</p>`
			} else {
				// 验证码验证
				sessionVcode, ok := utils.GlobalSessions.GetSessionData(r, "vcode")
				if !ok || strings.ToLower(vcode) != strings.ToLower(sessionVcode.(string)) {
					msg = `<p style="color:red;">验证码错误，请重新输入</p>`
				} else if username == "admin" && password == "123456" {
					msg = `<p style="color:green;">登录成功！</p>`
				} else {
					msg = `<p style="color:red;">用户名或密码错误</p>`
				}
			}
		}

		hiddenToken := `<input type="hidden" name="token" value="` + sessionToken + `" />`
		data := templates.NewPageData2(1, 6, msg+hiddenToken)
		renderer.RenderPage(w, "burteforce/bf_token.html", data)
	}
}

// 生成/读取 token（模拟 session 机制）
func getOrSetToken(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie("csrf_token")
	if err == nil {
		return cookie.Value
	}
	// 生成新 token
	rand.Seed(time.Now().UnixNano())
	token := strconv.FormatInt(rand.Int63(), 16)
	http.SetCookie(w, &http.Cookie{
		Name:  "csrf_token",
		Value: token,
		Path:  "/",
	})
	return token
}
