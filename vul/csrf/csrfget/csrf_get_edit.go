package csrfget

import (
	"net/http"
)

// CsrfGetEditHandler 用 GET 修改 cookie 中昵称
func CsrfGetEditHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nickname := r.URL.Query().Get("nickname")
		http.SetCookie(w, &http.Cookie{
			Name:  "csrf_get_nick",
			Value: nickname,
			Path:  "/",
		})
		http.Redirect(w, r, "/vul/csrf/csrfget/csrf_get", http.StatusFound)
	}
}
