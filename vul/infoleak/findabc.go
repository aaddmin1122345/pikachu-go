package infoleak

import (
	"net/http"
	"pikachu-go/templates"
)

// FindABCHandler 设置模拟登录 Cookie
func FindABCHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uname := r.URL.Query().Get("uname")
		pw := r.URL.Query().Get("pw")

		msg := "请输入用户名和密码作为 GET 参数，如：?uname=admin&pw=123456"

		if uname != "" && pw != "" {
			// 设置模拟敏感 cookie
			http.SetCookie(w, &http.Cookie{
				Name:  "abc[uname]",
				Value: uname,
				Path:  "/",
			})
			http.SetCookie(w, &http.Cookie{
				Name:  "abc[pw]",
				Value: pw,
				Path:  "/",
			})

			msg = "成功设置 cookie：abc[uname]=" + uname + " / abc[pw]=" + pw
		}

		data := templates.NewPageData2(66, 69, msg)
		renderer.RenderPage(w, "infoleak/findabc.html", data)
	}
}
