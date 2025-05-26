package xssblind

import (
	"crypto/md5"
	"crypto/sha1"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"pikachu-go/database"
	"pikachu-go/templates"
	"time"
)

// AdminLoginHandler 处理后台登录
func AdminLoginHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewPageData2(7, 14, "")

		htmlMsg := "<p>please input username and password!</p>"

		if r.Method == http.MethodPost {
			username := r.FormValue("username")
			password := r.FormValue("password")

			if username != "" && password != "" {
				// 获取数据库连接
				db := database.DB

				if db != nil {
					// 查询用户
					var id int
					var storedUsername string

					// 使用MD5进行密码加密后比较
					hashedPassword := fmt.Sprintf("%x", md5.Sum([]byte(password)))
					err := db.QueryRow("SELECT id, username FROM users WHERE username=$1 AND password=$2",
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

						// 重定向到管理页面
						http.Redirect(w, r, "/vul/xss/xssblind/admin", http.StatusFound)
						return
					} else if err != sql.ErrNoRows {
						// 数据库错误
						htmlMsg = "<p>Database error occurred!</p>"
					} else {
						// 用户名或密码错误
						htmlMsg = "<p>username or password error!</p>"
					}
				} else {
					htmlMsg = "<p>Database connection error!</p>"
				}
			} else {
				htmlMsg = "<p>please input username and password!</p>"
			}
		}

		data.HtmlMsg = template.HTML(htmlMsg)
		renderer.RenderPage(w, "xss/xssblind/admin_login.html", data)
	}
}
