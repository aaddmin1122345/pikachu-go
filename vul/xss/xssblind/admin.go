package xssblind

import (
	"fmt"
	"html/template"
	"net/http"
	"pikachu-go/database"
	"pikachu-go/templates"
	"strconv"
)

func AdminHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewPageData2(7, 14, "")

		// 检查登录状态
		_, err := r.Cookie("ant_uname")
		if err != nil {
			http.Redirect(w, r, "/vul/xss/xssblind/admin_login", http.StatusFound)
			return
		}

		// 处理退出登录
		if r.URL.Query().Get("logout") == "1" {
			http.SetCookie(w, &http.Cookie{Name: "ant_uname", Value: "", MaxAge: -1})
			http.SetCookie(w, &http.Cookie{Name: "ant_pw", Value: "", MaxAge: -1})
			http.Redirect(w, r, "/vul/xss/xssblind/admin_login", http.StatusFound)
			return
		}

		// 处理删除操作
		if id := r.URL.Query().Get("id"); id != "" {
			if idNum, err := strconv.Atoi(id); err == nil {
				db := database.DB
				if db != nil {
					_, err := db.Exec("DELETE FROM xssblind WHERE id=$1", idNum)
					if err != nil {
						// 可以在这里添加错误处理代码
					}
				}
			}
		}

		// 显示留言列表
		db := database.DB
		htmlContent := "你已经成功登录留言板后台,<a href=\"/vul/xss/xssblind/admin?logout=1\">退出登陆</a>"
		htmlContent += "<h2>用户反馈的意见列表：</h2>"
		htmlContent += "<table class=\"table table-bordered table-striped\">"
		htmlContent += "<tr><td>编号</td><td>时间</td><td>内容</td><td>姓名</td><td>操作</td></tr>"

		if db != nil {
			// 修改排序，让最新留言在前面
			rows, err := db.Query("SELECT id, time, content, name FROM xssblind ORDER BY time DESC")
			if err == nil {
				defer rows.Close()

				for rows.Next() {
					var id int
					var time, content, name string
					if err := rows.Scan(&id, &time, &content, &name); err == nil {
						// 这里直接输出用户输入内容，不进行过滤，会造成XSS漏洞
						htmlContent += fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%s</td><td>%s</td><td><a href=\"/vul/xss/xssblind/admin?id=%d\">删除</a></td></tr>",
							id, time, content, name, id)
					}
				}
			} else {
				htmlContent += fmt.Sprintf("<tr><td colspan='5'>获取留言列表失败: %s</td></tr>", err.Error())
			}
		} else {
			htmlContent += "<tr><td colspan='5'>数据库连接失败</td></tr>"
		}

		htmlContent += "</table>"

		data.HtmlMsg = template.HTML(htmlContent)
		renderer.RenderPage(w, "xss/xssblind/admin.html", data)
	}
}
