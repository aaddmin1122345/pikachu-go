package csrftoken

import (
	"net/http"
)

// TokenGetEditHandler 验证 token 后写入 nickname
func TokenGetEditHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nickname := r.URL.Query().Get("nickname")
		token := r.URL.Query().Get("token")

		cookie, err := r.Cookie("token_get_token")
		if err != nil || token != cookie.Value {
			http.Error(w, "非法请求，token 验证失败", http.StatusForbidden)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "token_get_nick",
			Value: nickname,
			Path:  "/",
		})
		http.Redirect(w, r, "/vul/csrf/csrftoken/token_get", http.StatusFound)
	}
}
