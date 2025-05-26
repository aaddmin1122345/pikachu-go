package csrfpost

import (
	"net/http"
	"pikachu-go/templates"
	"pikachu-go/utils"
)

// CsrfPostLoginHandler 处理CSRF POST登录
func CsrfPostLoginHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 检查是否已登录，如果已登录直接重定向到会员中心
		loggedIn, _ := utils.CheckCSRFLogin(r)
		if loggedIn {
			http.Redirect(w, r, "/vul/csrf/csrfpost/csrf_post", http.StatusFound)
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
					// 登录成功，重定向到会员中心
					http.Redirect(w, r, "/vul/csrf/csrfpost/csrf_post", http.StatusFound)
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
				
				<form method="post" action="/vul/csrf/csrfpost/csrf_post_login">
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
			</div>
		</div>
		`

		data := templates.NewPageData2(25, 28, html)
		renderer.RenderPage(w, "csrf/csrfpost/csrf_post_login.html", data)
	}
}
