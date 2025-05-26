package csrfget

import (
	"net/http"
	"pikachu-go/database"
	"pikachu-go/templates"
	"pikachu-go/utils"
)

// CsrfGetHandler 显示用户信息页面
func CsrfGetHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 检查用户是否登录
		loggedIn, username := utils.CheckCSRFLogin(r)
		if !loggedIn {
			http.Redirect(w, r, "/vul/csrf/login?module=get&redirect=/vul/csrf/csrfget/csrf_get", http.StatusFound)
			return
		}

		// 处理登出请求
		if r.URL.Query().Get("logout") == "1" {
			utils.ClearCSRFSession(w)
			http.Redirect(w, r, "/vul/csrf/login?module=get", http.StatusFound)
			return
		}

		// 从数据库查询用户信息
		var sex, phonenum, address, email string
		err := database.DB.QueryRow("SELECT sex, phonenum, address, email FROM member WHERE username = $1", username).Scan(&sex, &phonenum, &address, &email)
		if err != nil {
			http.Error(w, "无法获取用户信息", http.StatusInternalServerError)
			return
		}

		// 构建HTML内容
		html := `
		<div id="per_info">
		   <h1 class="per_title">hello,` + username + `,欢迎来到个人会员中心 | <a style="color:blue;" href="/vul/csrf/csrfget/csrf_get?logout=1">退出登录</a></h1>
		   <p class="per_name">姓名:` + username + `</p>
		   <p class="per_sex">性别:` + sex + `</p>
		   <p class="per_phone">手机:` + phonenum + `</p>
		   <p class="per_add">住址:` + address + `</p>
		   <p class="per_email">邮箱:` + email + `</p>
		   <a class="edit" href="/vul/csrf/csrfget/csrf_get_edit">修改个人信息</a>
		</div>
		`

		data := templates.NewPageData2(25, 27, html)
		renderer.RenderPage(w, "csrf/csrfget/csrf_get.html", data)
	}
}
