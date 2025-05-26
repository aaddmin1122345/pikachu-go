package csrftoken

import (
	"net/http"
	"pikachu-go/templates"
	"pikachu-go/utils"
)

// TokenGetLoginHandler 处理带Token保护的登录
func TokenGetLoginHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 检查是否已登录，如果已登录直接重定向到会员中心
		loggedIn, _ := utils.CheckCSRFLogin(r)
		if loggedIn {
			http.Redirect(w, r, "/vul/csrf/csrftoken/token_get", http.StatusFound)
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
					// 生成CSRF token并设置到cookie
					token := utils.GenerateCSRFToken()
					http.SetCookie(w, &http.Cookie{
						Name:     "csrf_token",
						Value:    token,
						Path:     "/",
						HttpOnly: true,
						MaxAge:   3600, // 1小时过期
					})

					// 登录成功，重定向到会员中心
					http.Redirect(w, r, "/vul/csrf/csrftoken/token_get", http.StatusFound)
					return
				} else {
					message = "登录失败，用户名或密码错误"
				}
			} else {
				message = "请输入用户名和密码"
			}
		}

		// 构建登录表单HTML
		html := `
		<div class="xss_form">
			<div class="xss_form_main">
				<h4 class="header blue lighter bigger">
					<i class="ace-icon fa fa-coffee green"></i>
					请输入登录信息
				</h4>
				
				<form method="post" action="/vul/csrf/csrftoken/token_get_login">
					<label>
						<span>
							<input type="text" name="username" placeholder="用户名" />
							<i class="ace-icon fa fa-user"></i>
						</span>
					</label>
					</br>
					
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
			html += "<p>" + message + "</p>"
		}

		html += `
				<p style="color:red">提示: 用户名有 vince/allen/kobe/grady/kevin/lucy/lili，密码全部是 123456</p>
				<p>此功能区域使用了CSRF Token来防护CSRF攻击</p>
			</div>
		</div>
		`

		data := templates.NewPageData2(25, 29, html)
		renderer.RenderPage(w, "csrf/csrftoken/token_get_login.html", data)
	}
}
