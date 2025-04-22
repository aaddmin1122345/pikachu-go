package sqliiu

import (
	"net/http"
	"pikachu-go/templates"
)

// SqliMemHandler 登录成功页面
func SqliMemHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("login_user")
		msg := "未登录，请先登录"

		if err == nil && cookie.Value != "" {
			msg = "欢迎你，" + cookie.Value
		}

		data := templates.NewPageData2(35, 41, msg)
		renderer.RenderPage(w, "sqli/sqli_iu/sqli_mem.html", data)
	}
}
