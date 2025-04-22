package burteforce

import (
	"net/http"
	"pikachu-go/templates"
	"strconv"
)

// BfServerHandler 服务端控制爆破尝试次数（简单 session 替代）
func BfServerHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tryCount int
		cookie, err := r.Cookie("bf_try")
		if err == nil {
			tryCount, _ = strconv.Atoi(cookie.Value)
		}

		msg := ""

		if r.Method == http.MethodPost {
			if tryCount >= 3 {
				msg = `<p style="color:red;">尝试次数过多，请稍后再试！</p>`
			} else {
				username := r.FormValue("username")
				password := r.FormValue("password")

				if username == "admin" && password == "123456" {
					msg = `<p style="color:green;">登录成功！</p>`
					tryCount = 0 // 登录成功则重置
				} else {
					tryCount++
					msg = `<p style="color:red;">用户名或密码错误（当前尝试次数：` + strconv.Itoa(tryCount) + `）</p>`
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
