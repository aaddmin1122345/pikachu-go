package xsspost

import (
	"crypto/md5"
	"crypto/sha1"
	"database/sql"
	"fmt"
	"net/http"
	"pikachu-go/db"
	"pikachu-go/templates"
	"time"
)

// PostLoginHandler 处理POST登录请求
func PostLoginHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewPageData2(7, 10, "")

		htmlMsg := "<p>please input username and password!</p>"

		if r.Method == http.MethodPost {
			username := r.FormValue("username")
			password := r.FormValue("password")

			if username != "" && password != "" {
				// 获取数据库连接
				database := db.GetDB()

				// 查询用户
				var id int
				var storedUsername string

				// 使用MD5进行密码加密后比较
				hashedPassword := fmt.Sprintf("%x", md5.Sum([]byte(password)))
				err := database.QueryRow("SELECT id, username FROM users WHERE username=? AND password=?",
					username, hashedPassword).Scan(&id, &storedUsername)

				if err == nil {
					// 登录成功，设置cookie
					expiration := time.Now().Add(1 * time.Hour)
					http.SetCookie(w, &http.Cookie{
						Name:    "ant_uname",
						Value:   username,
						Expires: expiration,
					})

					// 使用sha1(md5(password))存储密码cookie，与原版保持一致
					pwMd5 := md5.Sum([]byte(password))
					pwSha1 := sha1.Sum([]byte(fmt.Sprintf("%x", pwMd5)))
					http.SetCookie(w, &http.Cookie{
						Name:    "ant_pw",
						Value:   fmt.Sprintf("%x", pwSha1),
						Expires: expiration,
					})

					// 重定向到XSS反射页面
					http.Redirect(w, r, "/vul/xss/xsspost/xss_reflected_post", http.StatusFound)
					return
				} else if err != sql.ErrNoRows {
					// 数据库错误
					htmlMsg = "<p>Database error occurred!</p>"
				} else {
					// 用户名或密码错误
					htmlMsg = "<p>username or password error!</p>"
				}
			} else {
				htmlMsg = "<p>please input username and password!</p>"
			}
		}

		data.HtmlMsg = templates.HTML(htmlMsg)
		renderer.RenderPage(w, "xss/xsspost/post_login.html", data)
	}
}
