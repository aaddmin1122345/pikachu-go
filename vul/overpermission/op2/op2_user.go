package op2

import (
	"fmt"
	"net/http"
	"pikachu-go/templates"
)

// Op2UserHandler 处理用户中心页面
func Op2UserHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 检查是否登录
		loginStatus, sessionMap := checkOp2Login(r)
		if !loginStatus {
			// 未登录，重定向到登录页
			http.Redirect(w, r, "/vul/overpermission/op2/op2_login", http.StatusFound)
			return
		}

		username := sessionMap["username"].(string)
		isAdmin, _ := sessionMap["isAdmin"].(bool)

		// 如果请求包含退出登录参数
		if r.URL.Query().Get("logout") == "1" {
			// 清除会话
			clearOp2Session(w)
			http.Redirect(w, r, "/vul/overpermission/op2/op2_login", http.StatusFound)
			return
		}

		// 构建用户中心页面
		html := fmt.Sprintf(`
		<div class="user-center">
			<h2>用户中心</h2>
			<p>欢迎，%s！您的用户角色：%s</p>
			<p><a href="/vul/overpermission/op2/op2_user?logout=1">退出登录</a></p>
			
			<div class="menu">
				<h3>功能菜单</h3>
				<ul>
					<li><a href="#">个人资料</a></li>
					<li><a href="#">修改密码</a></li>
					<li><a href="#">我的消息</a></li>
					%s
				</ul>
			</div>
		</div>
		`, username, getRoleName(isAdmin), getAdminMenu(isAdmin))

		data := templates.NewPageData2(73, 78, html)
		data.PikaRoot = "/"
		renderer.RenderPage(w, "overpermission/op2/op2_user.html", data)
	}
}

// 获取角色名称
func getRoleName(isAdmin bool) string {
	if isAdmin {
		return "管理员"
	}
	return "普通用户"
}

// 获取管理员菜单项
func getAdminMenu(isAdmin bool) string {
	if isAdmin {
		return `<li><a href="/vul/overpermission/op2/op2_admin">管理员功能</a></li>`
	}
	return ""
}

// 清除会话
func clearOp2Session(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   "PIKASESSION",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}
