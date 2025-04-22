package op1

import (
	"database/sql"
	"fmt"
	"net/http"
	"pikachu-go/database"
	"pikachu-go/templates"
	"pikachu-go/utils"
)

// Op1MemHandler 处理会员信息页面，展示越权漏洞
func Op1MemHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 检查是否登录
		loginStatus, username := checkOp1Login(r)
		if !loginStatus {
			// 未登录，重定向到登录页
			http.Redirect(w, r, "/vul/overpermission/op1/op1_login", http.StatusFound)
			return
		}

		// 获取要查询的用户名（可能来自URL参数或会话中）
		queryUsername := r.URL.Query().Get("username")
		if queryUsername == "" {
			// 如果URL参数中没有指定，则使用当前登录用户
			queryUsername = username
		}

		html := ""
		// 默认情况下显示表单
		if r.Method == http.MethodGet && r.URL.Query().Get("submit") != "" {
			// 查询用户信息 - 这里存在水平越权漏洞，未验证查询的用户是否就是当前登录用户
			userInfo, err := getUserInfo(queryUsername)
			if err != nil {
				html = fmt.Sprintf("<p class='notice'>获取用户信息失败: %s</p>", err.Error())
			} else {
				// 展示用户信息
				html = fmt.Sprintf(`
				<div id="per_info">
					<h1 class="per_title">hello, %s, 你的具体信息如下：</h1>
					<p class="per_name">姓名: %s</p>
					<p class="per_sex">性别: %s</p>
					<p class="per_phone">手机: %s</p>
					<p class="per_add">住址: %s</p>
					<p class="per_email">邮箱: %s</p>
				</div>
				`, userInfo.Username, userInfo.Username, userInfo.Sex, userInfo.PhoneNum, userInfo.Address, userInfo.Email)
			}
		}

		// 如果请求包含退出登录参数
		if r.URL.Query().Get("logout") == "1" {
			// 清除会话
			clearOp1Session(w)
			http.Redirect(w, r, "/vul/overpermission/op1/op1_login", http.StatusFound)
			return
		}

		// 渲染页面，添加一个查询自己信息的表单和当前登录用户显示
		html = `
		<p class="mem_title">欢迎来到个人信息中心，当前用户: ` + username + ` | <a style="color:blue;" href="/vul/overpermission/op1/op1_mem?logout=1">退出登录</a></p>
		<form class="msg1" method="get">
			<input type="hidden" name="username" value="` + username + `" />
			<input type="submit" name="submit" value="点击查看个人信息" />
		</form>
		` + html

		data := templates.NewPageData2(73, 77, html)
		data.PikaRoot = "/"
		renderer.RenderPage(w, "overpermission/op1/op1_mem.html", data)
	}
}

// 用户信息结构体
type UserInfo struct {
	Username string
	Sex      string
	PhoneNum string
	Address  string
	Email    string
}

// 获取用户信息
func getUserInfo(username string) (*UserInfo, error) {
	db := database.DB
	if db == nil {
		return nil, fmt.Errorf("数据库连接失败")
	}

	// 查询用户信息
	row := db.QueryRow("SELECT username, sex, phonenum, address, email FROM member WHERE username = ?", username)

	var user UserInfo
	err := row.Scan(&user.Username, &user.Sex, &user.PhoneNum, &user.Address, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, err
	}

	return &user, nil
}

// 检查是否登录
func checkOp1Login(r *http.Request) (bool, string) {
	sessionData, ok := utils.GlobalSessions.GetSessionData(r, "op1")
	if !ok {
		return false, ""
	}

	sessionMap, ok := sessionData.(map[string]interface{})
	if !ok {
		return false, ""
	}

	username, ok := sessionMap["username"].(string)
	if !ok {
		return false, ""
	}

	return true, username
}

// 清除会话
func clearOp1Session(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   "PIKASESSION",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}
