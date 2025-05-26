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

				// 添加水平越权提示
				if queryUsername != username {
					html += `
					<div style="margin-top:20px;padding:10px;background-color:#ffdddd;border-left:6px solid #f44336;">
						<strong>安全提示!</strong> 您刚刚查看了其他用户的信息 (水平越权漏洞)。<br>
						这是因为后端仅验证了登录状态，但没有验证当前登录用户是否有权查看所请求的用户信息。<br>
						正确的做法是：验证请求的用户名与当前登录用户是否一致。
					</div>
					`
				}
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
		welcomeMsg := fmt.Sprintf(`<p class="mem_title">欢迎来到个人信息中心，当前用户: %s | <a style="color:blue;" href="/vul/overpermission/op1/op1_mem?logout=1">退出登录</a></p>`, username)

		// 添加正常查询表单和越权测试表单
		formHtml := fmt.Sprintf(`
		<form class="msg1" method="get" style="margin-bottom:10px;">
			<input type="hidden" name="username" value="%s" />
			<input type="submit" name="submit" value="点击查看个人信息" />
		</form>
		
		<div style="margin:15px 0;border-top:1px dashed #ccc;padding-top:15px;">
			<h4>水平越权漏洞测试</h4>
			<p>尝试查看其他用户信息 (可用测试账号: pikachu, admin, lucy, lili, kobe)</p>
			<form class="msg1" method="get">
				<input type="text" name="username" placeholder="输入要查询的用户名" style="width:200px;margin-right:10px;" />
				<input type="submit" name="submit" value="查看此用户信息" />
			</form>
			<div style="margin-top:10px;padding:8px;background-color:#e7f3fe;border-left:6px solid #2196F3;">
				<strong>漏洞说明:</strong> 此页面演示了水平越权漏洞。登录用户可以查看任意其他用户的个人信息，而不仅限于自己的信息。
			</div>
		</div>
		`, username)

		finalHtml := welcomeMsg + formHtml + html

		data := templates.NewPageData2(73, 77, finalHtml)
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
