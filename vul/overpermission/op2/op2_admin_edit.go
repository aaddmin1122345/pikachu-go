package op2

import (
	"fmt"
	"net/http"
	"pikachu-go/database"
	"pikachu-go/templates"
	"strconv"
)

// Op2AdminEditHandler 处理管理员编辑用户功能页面（存在垂直越权漏洞）
func Op2AdminEditHandler(renderer templates.Renderer) http.HandlerFunc {
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

		// 获取要编辑的用户ID
		userID := r.URL.Query().Get("id")
		if userID == "" {
			http.Redirect(w, r, "/vul/overpermission/op2/op2_admin", http.StatusFound)
			return
		}

		msg := ""
		// 处理表单提交
		if r.Method == http.MethodPost {
			newUsername := r.FormValue("username")
			newRole := r.FormValue("role")

			if newUsername == "" {
				msg = "<p style='color:red;'>用户名不能为空</p>"
			} else {
				// 更新用户信息
				if updateUser(userID, newUsername, newRole) {
					msg = "<p style='color:green;'>用户信息更新成功</p>"
				} else {
					msg = "<p style='color:red;'>用户信息更新失败</p>"
				}
			}
		}

		// 获取用户信息
		userInfo, err := getUserByID(userID)
		if err != nil {
			http.Error(w, "获取用户信息失败: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 构建编辑表单
		html := fmt.Sprintf(`
		<div class="admin-edit">
			<h2>编辑用户</h2>
			<p>当前登录用户: %s</p>
			<p><a href="/vul/overpermission/op2/op2_admin">返回用户列表</a></p>
			
			<form method="post">
				<div class="form-group">
					<label>用户ID:</label>
					<input type="text" value="%s" disabled />
				</div>
				<div class="form-group">
					<label>用户名:</label>
					<input type="text" name="username" value="%s" required />
				</div>
				<div class="form-group">
					<label>角色:</label>
					<select name="role">
						<option value="user" %s>普通用户</option>
						<option value="admin" %s>管理员</option>
					</select>
				</div>
				<div class="form-group">
					<input type="submit" value="保存修改" />
				</div>
			</form>
			%s
		</div>
		`, username, userID, userInfo.Username,
			getSelected(userInfo.Role, "user"), getSelected(userInfo.Role, "admin"),
			msg)

		data := templates.NewPageData2(73, 80, html)
		data.PikaRoot = "/"
		renderer.RenderPage(w, "overpermission/op2/op2_admin_edit.html", data)
	}
}

// 用户信息结构
type UserDetail struct {
	ID       int
	Username string
	Role     string
}

// 通过ID获取用户信息
func getUserByID(userID string) (*UserDetail, error) {
	db := database.DB
	if db == nil {
		return nil, fmt.Errorf("数据库连接失败")
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		return nil, fmt.Errorf("无效的用户ID")
	}

	var user UserDetail
	err = db.QueryRow("SELECT id, username, role FROM users WHERE id = ?", id).Scan(&user.ID, &user.Username, &user.Role)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %s", err.Error())
	}

	return &user, nil
}

// 更新用户信息
func updateUser(userID, username, role string) bool {
	db := database.DB
	if db == nil {
		return false
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		return false
	}

	_, err = db.Exec("UPDATE users SET username = ?, role = ? WHERE id = ?", username, role, id)
	return err == nil
}

// 获取下拉选项的selected属性
func getSelected(current, option string) string {
	if current == option {
		return "selected"
	}
	return ""
}
