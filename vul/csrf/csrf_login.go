package csrf

import (
	"html/template"
	"net/http"
	"pikachu-go/templates"
	"pikachu-go/utils"
)

// CsrfLoginHandler 统一的CSRF登录处理函数
func CsrfLoginHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取目标重定向URL
		redirectURL := r.URL.Query().Get("redirect")
		if redirectURL == "" {
			// 默认重定向到GET版本
			redirectURL = "/vul/csrf/csrfget/csrf_get"
		}

		// 获取目标模块（GET、POST或Token）
		module := r.URL.Query().Get("module")
		if module == "" {
			module = "get" // 默认是GET
		}

		// 检查是否已登录，如果已登录直接重定向
		loggedIn, _ := utils.CheckCSRFLogin(r)
		if loggedIn {
			http.Redirect(w, r, redirectURL, http.StatusFound)
			return
		}

		// 处理登录逻辑
		message := ""
		if r.Method == http.MethodPost {
			username := r.FormValue("username")
			password := r.FormValue("password")

			if username != "" && password != "" {
				// 设置会话
				if utils.SetCSRFSession(w, username, password) {
					// 如果是Token模块，额外设置CSRF token
					if module == "token" {
						token := utils.GenerateCSRFToken()
						http.SetCookie(w, &http.Cookie{
							Name:     "csrf_token",
							Value:    token,
							Path:     "/",
							HttpOnly: true,
							MaxAge:   3600, // 1小时过期
						})
					}

					// 登录成功，重定向到目标页面
					http.Redirect(w, r, redirectURL, http.StatusFound)
					return
				} else {
					message = "登录失败，用户名或密码错误"
				}
			} else {
				message = "请输入用户名和密码"
			}
		}

		// 构建登录表单HTML
		var pageTitle string
		switch module {
		case "post":
			pageTitle = "CSRF POST测试 - 登录"
		case "token":
			pageTitle = "CSRF Token测试 - 登录"
		default:
			pageTitle = "CSRF GET测试 - 登录"
		}

		html := `
		<div class="xss_form">
			<div class="xss_form_main">
				<h4 class="header blue lighter bigger">
					<i class="ace-icon fa fa-coffee green"></i>
					` + pageTitle + `
				</h4>
				
				<form method="post" action="/vul/csrf/login?module=` + module + `&redirect=` + redirectURL + `">
					<label>
						<span>
							<input type="text" name="username" placeholder="用户名" />
							<i class="ace-icon fa fa-user"></i>
						</span>
					</label>
					<br>
					
					<label>
						<span>
							<input type="password" name="password" placeholder="密码" />
							<i class="ace-icon fa fa-lock"></i>
						</span>
					</label>
					
					<div class="space"></div>
					
					<div class="clearfix">
						<label><input class="submit" name="submit" type="submit" value="登录" /></label>
					</div>
				</form>
		`

		if message != "" {
			html += "<p style='color:red'>" + message + "</p>"
		}

		html += `
				<p style="color:red">提示: 用户名有 vince/allen/kobe/grady/kevin/lucy/lili，密码全部是 123456</p>
			`

		if module == "token" {
			html += "<p>此功能区域使用了CSRF Token来防护CSRF攻击</p>"
		}

		html += `</div></div>`

		// 根据模块选择相应的页面数据
		var data templates.PageData
		activeIndex := 27 // 默认GET模块
		if module == "post" {
			activeIndex = 28
		} else if module == "token" {
			activeIndex = 29
		}

		active := make([]string, 130)
		active[25] = "active open"
		active[activeIndex] = "active"

		data = templates.PageData{
			Active:  active,
			HtmlMsg: template.HTML(html),
		}

		renderer.RenderPage(w, "csrf/csrf_login.html", data)
	}
}
