package infoleak

import (
	"fmt"
	"net/http"
	"pikachu-go/templates"
)

// FindABCHandler 设置模拟登录 Cookie
func FindABCHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uname := r.URL.Query().Get("uname")
		pw := r.URL.Query().Get("pw")

		msg := "请输入用户名和密码作为 GET 参数，例如：?uname=admin&pw=123456"

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

			msg = fmt.Sprintf(`
				<p>您输入的信息已提交！</p>
				<p style="color:#337ab7; font-weight:bold;">Cookie已设置：</p>
				<pre style="background:#f5f5f5; padding:10px; border:1px solid #ddd;">
abc[uname] = %s
abc[pw] = %s
				</pre>
				<p style="margin-top:15px;">提示：您现在可以查看源代码或检查其他信息...</p>
			`, uname, pw)
		}

		data := templates.NewPageData2(85, 87, msg)
		renderer.RenderPage(w, "infoleak/findabc.html", data)
	}
}
