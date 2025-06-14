package op2

import (
	"fmt"
	"net/http"
	"pikachu-go/database"
	"pikachu-go/templates"
)

// Op2AdminHandler 处理管理员功能页面（存在垂直越权漏洞）
func Op2AdminHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 检查是否登录
		loginStatus, sessionMap := checkOp2Login(r)
		if !loginStatus {
			// 未登录，重定向到登录页
			http.Redirect(w, r, "/vul/overpermission/op2/op2_login", http.StatusFound)
			return
		}

		username := sessionMap["username"].(string)

		// 注意：这里存在垂直越权漏洞，没有验证用户是否为管理员
		// 正确的做法是：
		// isAdmin, _ := sessionMap["isAdmin"].(bool)
		// if !isAdmin {
		//     http.Error(w, "权限不足", http.StatusForbidden)
		//     return
		// }

		html := fmt.Sprintf(`
		<div class="admin-panel">
			<h2>管理员控制面板</h2>
			<p>当前登录用户: %s</p>
			<p><a href="/vul/overpermission/op2/op2_user?logout=1">退出登录</a></p>
			
			<div style="margin:15px 0;padding:10px;background-color:#ffdddd;border-left:6px solid #f44336;">
				<strong>安全提示!</strong> 您正在访问仅限管理员的页面 (垂直越权漏洞)。<br>
				这是因为后端仅验证了登录状态，但没有验证用户是否具有管理员权限。<br>
				正确的做法是：验证当前用户是否拥有管理员角色。
			</div>
			
			<div class="user-list">
				<h3>用户列表</h3>
				<table border="1" cellpadding="5">
					<tr>
						<th>ID</th>
						<th>用户名</th>
						<th>角色</th>
						<th>操作</th>
					</tr>
					%s
				</table>
			</div>
			
			<div style="margin-top:15px;padding:8px;background-color:#e7f3fe;border-left:6px solid #2196F3;">
				<strong>漏洞说明:</strong> 此页面演示了垂直越权漏洞。普通用户可以访问仅限管理员的功能页面，因为系统没有进行权限检查。
			</div>
		</div>
		`, username, getUserListHTML())

		data := templates.NewPageData2(73, 79, html)
		data.PikaRoot = "/"
		renderer.RenderPage(w, "overpermission/op2/op2_admin.html", data)
	}
}

// 获取用户列表HTML
func getUserListHTML() string {
	db := database.DB
	if db == nil {
		return "<tr><td colspan='4'>数据库连接失败</td></tr>"
	}

	rows, err := db.Query("SELECT id, username, role FROM users")
	if err != nil {
		return fmt.Sprintf("<tr><td colspan='4'>获取用户列表失败: %s</td></tr>", err.Error())
	}
	defer rows.Close()

	html := ""
	for rows.Next() {
		var id int
		var username, role string
		if err := rows.Scan(&id, &username, &role); err != nil {
			continue
		}

		html += fmt.Sprintf(`
		<tr>
			<td>%d</td>
			<td>%s</td>
			<td>%s</td>
			<td>
				<a href="/vul/overpermission/op2/op2_admin_edit?id=%d">编辑</a> | 
				<a href="/vul/overpermission/op2/op2_admin?delete=%d" onclick="return confirm('确定删除?')">删除</a>
			</td>
		</tr>
		`, id, username, role, id, id)
	}

	if html == "" {
		html = "<tr><td colspan='4'>没有用户</td></tr>"
	}

	return html
}
