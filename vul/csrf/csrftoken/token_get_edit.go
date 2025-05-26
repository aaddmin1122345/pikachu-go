package csrftoken

import (
	"log"
	"net/http"
	"pikachu-go/database"
	"pikachu-go/templates"
	"pikachu-go/utils"
)

// TokenGetEditHandler 带有CSRF Token保护的编辑功能
func TokenGetEditHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 检查用户是否登录
		loggedIn, username := utils.CheckCSRFLogin(r)
		if !loggedIn {
			http.Redirect(w, r, "/vul/csrf/login?module=token&redirect=/vul/csrf/csrftoken/token_get_edit", http.StatusFound)
			return
		}

		// 获取会话中的CSRF token
		session, _ := r.Cookie("csrf_token")
		sessionToken := ""
		if session != nil {
			sessionToken = session.Value
		}

		// 处理编辑提交
		message := ""
		if r.URL.Query().Get("submit") != "" {
			// 获取表单数据和token
			formToken := r.URL.Query().Get("token")
			sex := r.URL.Query().Get("sex")
			phonenum := r.URL.Query().Get("phonenum")
			address := r.URL.Query().Get("add")
			email := r.URL.Query().Get("email")

			// 检查必填字段和token
			if sex != "" && phonenum != "" && address != "" && email != "" {
				// 验证CSRF token
				if formToken == sessionToken && formToken != "" {
					// 更新数据库
					_, err := database.DB.Exec(
						"UPDATE member SET sex=$1, phonenum=$2, address=$3, email=$4 WHERE username=$5",
						sex, phonenum, address, email, username,
					)

					if err != nil {
						message = "修改失败，请重试"
						log.Printf("Token保护修改失败: %v", err)
					} else {
						// 修改成功，重定向到会员中心
						http.Redirect(w, r, "/vul/csrf/csrftoken/token_get", http.StatusFound)
						return
					}
				} else {
					message = "无效的安全令牌，可能是CSRF攻击尝试！"
				}
			}
		}

		// 获取当前用户信息用于表单显示
		var sex, phonenum, address, email string
		err := database.DB.QueryRow("SELECT sex, phonenum, address, email FROM member WHERE username = $1", username).Scan(&sex, &phonenum, &address, &email)
		if err != nil {
			http.Error(w, "无法获取用户信息", http.StatusInternalServerError)
			return
		}

		// 生成新的CSRF token并设置到cookie
		token := utils.GenerateCSRFToken()
		http.SetCookie(w, &http.Cookie{
			Name:     "csrf_token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   3600, // 1小时过期
		})

		// 构建编辑表单HTML - 包含隐藏的token字段
		html := `
		<div id="per_info">
		   <form method="get">
		   <h1 class="per_title">hello,` + username + `,欢迎来到个人会员中心 | <a style="color:blue;" href="/vul/csrf/csrftoken/token_get?logout=1">退出登录</a></h1>
		   <p class="per_name">姓名:` + username + `</p>
		   <p class="per_sex">性别:<input type="text" name="sex" value="` + sex + `"/></p>
		   <p class="per_phone">手机:<input class="phonenum" type="text" name="phonenum" value="` + phonenum + `"/></p>    
		   <p class="per_add">住址:<input class="add" type="text" name="add" value="` + address + `"/></p> 
		   <p class="per_email">邮箱:<input class="email" type="text" name="email" value="` + email + `"/></p>
		   <input type="hidden" name="token" value="` + token + `" /> 
		   <input class="sub" type="submit" name="submit" value="submit"/>
		   </form>
		</div>
		<div style="margin-top: 20px;">
			<p><strong>安全提示：</strong>此页面已添加CSRF Token保护机制，防止跨站请求伪造攻击。</p>
			<p>您的Token是: ` + token + `</p>
		</div>
		`
		if message != "" {
			html += "<p style='color:red'>" + message + "</p>"
		}

		data := templates.NewPageData2(25, 29, html)
		renderer.RenderPage(w, "csrf/csrftoken/token_get_edit.html", data)
	}
}
