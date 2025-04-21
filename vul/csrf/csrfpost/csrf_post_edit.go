package csrfpost

import (
	"net/http"
)

// CsrfPostEditHandler 接收 POST 请求修改 cookie
func CsrfPostEditHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			nick := r.FormValue("nickname")
			http.SetCookie(w, &http.Cookie{
				Name:  "csrf_post_nick",
				Value: nick,
				Path:  "/",
			})
		}
		http.Redirect(w, r, "/vul/csrf/csrfpost/csrf_post", http.StatusFound)
	}
}
