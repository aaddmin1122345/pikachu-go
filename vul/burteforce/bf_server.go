package burteforce

import (
	"net/http"
	"pikachu-go/templates"
	"pikachu-go/utils"
	"strconv"
	"strings"
)

// BfServerHandler 服务端控制爆破尝试次数（添加验证码验证）
func BfServerHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tryCount int
		cookie, err := r.Cookie("bf_try")
		if err == nil {
			tryCount, _ = strconv.Atoi(cookie.Value)
		}

		msg := ""

		if r.Method == http.MethodPost {
			username := r.FormValue("username")
			password := r.FormValue("password")
			vcode := r.FormValue("vcode")

			// 验证表单输入项是否为空
			if username == "" {
				msg = `<p style="color:red;">用户名不能为空</p>`
			} else if password == "" {
				msg = `<p style="color:red;">密码不能为空</p>`
			} else if vcode == "" {
				msg = `<p style="color:red;">验证码不能为空</p>`
			} else {
				// 验证码验证
				sessionVcode, ok := utils.GlobalSessions.GetSessionData(r, "vcode")
				if !ok || strings.ToLower(vcode) != strings.ToLower(sessionVcode.(string)) {
					msg = `<p style="color:red;">验证码错误，请重新输入</p>`
				} else {
					// 验证用户名和密码
					if username == "admin" && password == "123456" {
						msg = `<p style="color:green;">登录成功！</p>`
						tryCount = 0 // 登录成功则重置
					} else {
						tryCount++
						msg = `<p style="color:red;">用户名或密码错误（当前尝试次数：` + strconv.Itoa(tryCount) + `）</p>`
					}
				}
			}
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "bf_try",
			Value: strconv.Itoa(tryCount),
			Path:  "/",
		})

		data := templates.NewPageData2(1, 4, msg)
		renderer.RenderPage(w, "burteforce/bf_server.html", data)
	}
}
