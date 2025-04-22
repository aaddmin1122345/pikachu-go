package burteforce

import (
	"net/http"
	"pikachu-go/templates"
	"pikachu-go/utils"
	"strings"
)

// BfClientHandler 使用前端 JS 限制次数（服务端不判断）
func BfClientHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := ""
		if r.Method == http.MethodPost {
			username := r.FormValue("username")
			password := r.FormValue("password")
			vcode := r.FormValue("vcode")

			// 验证码验证仅作为示例，实际验证依然在前端JS中
			sessionVcode, ok := utils.GlobalSessions.GetSessionData(r, "vcode")
			if !ok || strings.ToLower(vcode) != strings.ToLower(sessionVcode.(string)) {
				msg = `<p style="color:red;">验证码错误</p>`
			} else if username == "admin" && password == "123456" {
				msg = `<p style="color:green;">登录成功！</p>`
			} else {
				msg = `<p style="color:red;">用户名或密码错误</p>`
			}
		}

		data := templates.NewPageData2(1, 5, msg)
		renderer.RenderPage(w, "burteforce/bf_client.html", data)
	}
}
