package sqliiu

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"pikachu-go/database"
	"pikachu-go/templates"
)

// SqliLoginHandler 登录页面
func SqliLoginHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := ""
		if r.Method == http.MethodPost {
			username := r.FormValue("username")
			password := r.FormValue("password")

			// 使用全局数据库连接
			db := database.DB
			if db == nil {
				http.Error(w, "数据库连接失败", http.StatusInternalServerError)
				return
			}

			// 使用MD5进行密码加密后比较
			hashedPassword := fmt.Sprintf("%x", md5.Sum([]byte(password)))

			// 故意保留SQL注入漏洞，直接拼接SQL
			query := fmt.Sprintf("SELECT id FROM member WHERE username = '%s' AND password = '%s'", username, hashedPassword)
			row := db.QueryRow(query)

			var id int
			if err := row.Scan(&id); err == nil {
				http.SetCookie(w, &http.Cookie{
					Name:  "login_user",
					Value: username,
					Path:  "/",
				})
				http.Redirect(w, r, "/vul/sqli/sqli_iu/sqli_mem", http.StatusFound)
				return
			}
			msg = "登录失败"
		}

		data := templates.NewPageData2(35, 41, msg)
		renderer.RenderPage(w, "sqli/sqli_iu/sqli_login.html", data)
	}
}
