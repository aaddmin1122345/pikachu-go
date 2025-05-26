package csrfget

import (
	"log"
	"net/http"
	"pikachu-go/database"
	"pikachu-go/templates"
	"pikachu-go/utils"
)

// CsrfGetEditHandler 通过GET请求修改用户信息 - 故意存在CSRF漏洞
func CsrfGetEditHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 检查用户是否登录
		loggedIn, username := utils.CheckCSRFLogin(r)
		if !loggedIn {
			http.Redirect(w, r, "/vul/csrf/login?module=get&redirect=/vul/csrf/csrfget/csrf_get_edit", http.StatusFound)
			return
		}

		// 处理编辑提交
		message := ""
		if r.URL.Query().Get("submit") != "" {
			// 获取表单数据
			sex := r.URL.Query().Get("sex")
			phonenum := r.URL.Query().Get("phonenum")
			address := r.URL.Query().Get("add")
			email := r.URL.Query().Get("email")

			// 检查必填字段
			if sex != "" && phonenum != "" && address != "" && email != "" {
				// 更新数据库 - 故意使用GET方法，造成CSRF漏洞
				_, err := database.DB.Exec(
					"UPDATE member SET sex=$1, phonenum=$2, address=$3, email=$4 WHERE username=$5",
					sex, phonenum, address, email, username,
				)

				if err != nil {
					message = "修改失败，请重试"
					log.Printf("CSRF GET修改失败: %v", err)
				} else {
					// 修改成功，重定向到会员中心
					http.Redirect(w, r, "/vul/csrf/csrfget/csrf_get", http.StatusFound)
					return
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

		// 构建编辑表单HTML
		html := `
		<div id="per_info">
		   <form method="get">
		   <h1 class="per_title">hello,` + username + `,欢迎来到个人会员中心 | <a style="color:blue;" href="/vul/csrf/csrfget/csrf_get?logout=1">退出登录</a></h1>
		   <p class="per_name">姓名:` + username + `</p>
		   <p class="per_sex">性别:<input type="text" name="sex" value="` + sex + `"/></p>
		   <p class="per_phone">手机:<input class="phonenum" type="text" name="phonenum" value="` + phonenum + `"/></p>    
		   <p class="per_add">住址:<input class="add" type="text" name="add" value="` + address + `"/></p> 
		   <p class="per_email">邮箱:<input class="email" type="text" name="email" value="` + email + `"/></p> 
		   <input class="sub" type="submit" name="submit" value="submit"/>
		   </form>
		</div>
		`
		if message != "" {
			html += "<p>" + message + "</p>"
		}

		data := templates.NewPageData2(25, 27, html)
		renderer.RenderPage(w, "csrf/csrfget/csrf_get_edit.html", data)
	}
}
