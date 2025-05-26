package op2

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"net/http"
	"pikachu-go/database"
	"pikachu-go/templates"
	"pikachu-go/utils"
)

// Op2LoginHandler 垂直越权登录页面
func Op2LoginHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := ""

		// 检查是否已经登录
		if loginStatus, _ := checkOp2Login(r); loginStatus {
			// 已登录，重定向到用户中心
			http.Redirect(w, r, "/vul/overpermission/op2/op2_user", http.StatusFound)
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
				isValid, isAdmin := validateOp2Login(username, password)
				if isValid {
					// 登录成功，设置会话
					sessionData := map[string]interface{}{
						"username": username,
						"password": password, // 实际应用中应存储加密后的密码或令牌
						"isAdmin":  isAdmin,
					}
					utils.GlobalSessions.SetSessionData(w, r, "op2", sessionData)

					// 重定向到用户中心
					http.Redirect(w, r, "/vul/overpermission/op2/op2_user", http.StatusFound)
					return
				} else {
					msg = "<p style='color:red;'>用户名或密码错误</p>"
				}
			}
		}

		data := templates.NewPageData2(73, 76, msg)
		data.PikaRoot = "/"
		renderer.RenderPage(w, "overpermission/op2/op2_login.html", data)
	}
}

// 检查是否登录
func checkOp2Login(r *http.Request) (bool, map[string]interface{}) {
	sessionData, ok := utils.GlobalSessions.GetSessionData(r, "op2")
	if !ok {
		return false, nil
	}

	sessionMap, ok := sessionData.(map[string]interface{})
	if !ok {
		return false, nil
	}

	_, ok = sessionMap["username"].(string)
	if !ok {
		return false, nil
	}

	return true, sessionMap
}

// 验证登录并检查是否为管理员
func validateOp2Login(username, password string) (bool, bool) {
	db := database.DB
	if db == nil {
		return false, false
	}

	// 查询用户
	var dbPassword string
	var role string
	err := db.QueryRow("SELECT password, role FROM users WHERE username = $1", username).Scan(&dbPassword, &role)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, false // 用户不存在
		}
		return false, false // 数据库错误
	}

	// 使用MD5进行密码加密后比较
	hashedPassword := fmt.Sprintf("%x", md5.Sum([]byte(password)))
	if hashedPassword != dbPassword {
		return false, false
	}

	// 检查是否为管理员
	isAdmin := (role == "admin")
	return true, isAdmin
}
