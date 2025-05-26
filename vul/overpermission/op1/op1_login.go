package op1

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"net/http"
	"pikachu-go/database"
	"pikachu-go/templates"
	"pikachu-go/utils"
)

// Op1LoginHandler 水平越权登录页面
func Op1LoginHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := ""

		// 检查是否已经登录
		if loginStatus, _ := checkOp1Login(r); loginStatus {
			// 已登录，重定向到会员中心
			http.Redirect(w, r, "/vul/overpermission/op1/op1_mem", http.StatusFound)
			return
		}

		// 处理登录请求
		if r.Method == http.MethodPost {
			username := r.FormValue("username")
			password := r.FormValue("password")

			// 验证登录
			if username == "" || password == "" {
				msg = "<p style='color:red;'>用户名和密码不能为空</p>"
			} else {
				// 验证用户名和密码
				if validateLogin(username, password) {
					// 登录成功，设置会话
					sessionData := map[string]interface{}{
						"username": username,
						"password": password, // 实际应用中应存储加密后的密码或令牌
					}
					utils.GlobalSessions.SetSessionData(w, r, "op1", sessionData)

					// 重定向到会员中心
					http.Redirect(w, r, "/vul/overpermission/op1/op1_mem", http.StatusFound)
					return
				} else {
					msg = "<p style='color:red;'>用户名或密码错误</p>"
				}
			}
		}

		// 添加测试账号提示
		if msg == "" {
			msg = `<div class="alert alert-info" style="margin-top: 20px;">
                本页面仅为演示，可用的测试账号：<br>
                用户名：pikachu，密码：123456<br>
                用户名：admin，密码：123456<br>
                用户名：lucy，密码：123456<br>
                用户名：lili，密码：123456<br>
                用户名：kobe，密码：123456
            </div>`
		}

		data := templates.NewPageData2(73, 75, msg)
		data.PikaRoot = "/"
		renderer.RenderPage(w, "overpermission/op1/op1_login.html", data)
	}
}

// 验证登录
func validateLogin(username, password string) bool {
	db := database.DB
	if db == nil {
		return false
	}

	// 查询用户
	var dbPassword string
	err := db.QueryRow("SELECT pw FROM member WHERE username = $1", username).Scan(&dbPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false // 用户不存在
		}
		return false // 数据库错误
	}

	// 使用MD5进行密码加密后比较
	hashedPassword := fmt.Sprintf("%x", md5.Sum([]byte(password)))
	return hashedPassword == dbPassword
}
