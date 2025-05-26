package csrf

import (
	"fmt"
	"html/template"
	"net/http"
	"pikachu-go/database"
	"pikachu-go/templates"
	"pikachu-go/utils"
)

// DebugLoginHandler 简单的调试登录处理函数
func DebugLoginHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		message := ""
		status := ""

		// 检查数据库
		dbStatus := "数据库状态：正常"
		rows, err := queryMembers()
		if err != nil {
			dbStatus = fmt.Sprintf("数据库错误：%v", err)
		}

		// 处理登录请求
		if r.Method == http.MethodPost {
			username := r.FormValue("username")
			password := r.FormValue("password")

			if username != "" && password != "" {
				if utils.SetCSRFSession(w, username, password) {
					message = "登录成功！"
					status = "success"
				} else {
					message = fmt.Sprintf("登录失败！用户名：%s，密码：%s", username, password)
					status = "error"
				}
			} else {
				message = "请输入用户名和密码"
				status = "warning"
			}
		}

		// 构建HTML
		html := `
		<div style="margin: 20px auto; max-width: 600px; padding: 20px; border: 1px solid #ccc; border-radius: 5px;">
			<h2>CSRF 调试登录页面</h2>
			<div style="background-color: #f5f5f5; padding: 10px; margin-bottom: 20px; border-radius: 5px;">
				` + dbStatus + `
				<div>
					<h4>数据库用户列表：</h4>
					<pre>` + rows + `</pre>
				</div>
			</div>
		`

		if message != "" {
			var color string
			switch status {
			case "success":
				color = "green"
			case "error":
				color = "red"
			default:
				color = "orange"
			}
			html += `<div style="color: ` + color + `; margin-bottom: 15px; padding: 10px; border: 1px solid ` + color + `; border-radius: 5px;">` + message + `</div>`
		}

		html += `
			<form method="post" action="/vul/csrf/debug_login" style="margin-top: 20px;">
				<div style="margin-bottom: 15px;">
					<label for="username">用户名：</label>
					<input type="text" id="username" name="username" style="padding: 5px; width: 100%;">
				</div>
				<div style="margin-bottom: 15px;">
					<label for="password">密码：</label>
					<input type="password" id="password" name="password" style="padding: 5px; width: 100%;">
				</div>
				<button type="submit" style="padding: 8px 15px; background-color: #4CAF50; color: white; border: none; border-radius: 4px; cursor: pointer;">登录</button>
			</form>
			
			<div style="margin-top: 20px; font-size: 0.9em; color: #666;">
				<p>提示: 用户名有 vince/allen/kobe/grady/kevin/lucy/lili，密码全部是 123456</p>
			</div>
		</div>
		`

		// 渲染页面
		data := templates.PageData{
			Active:  make([]string, 130),
			HtmlMsg: template.HTML(html),
		}

		renderer.RenderPage(w, "csrf/csrf_login.html", data)
	}
}

// 查询member表中的用户
func queryMembers() (string, error) {
	rows, err := database.DB.Query("SELECT id, username, pw, sex FROM member LIMIT 10")
	if err != nil {
		return "查询失败", err
	}
	defer rows.Close()

	result := ""
	for rows.Next() {
		var id int
		var username, pw, sex string
		if err := rows.Scan(&id, &username, &pw, &sex); err != nil {
			return "数据读取失败", err
		}
		result += fmt.Sprintf("ID: %d, 用户名: %s, 密码哈希: %s, 性别: %s\n", id, username, pw, sex)
	}

	if result == "" {
		return "数据库表为空", nil
	}

	return result, nil
}
