package xss

import (
	"fmt"
	"html/template"
	"net/http"
	"pikachu-go/database"
	"pikachu-go/templates"
	"strconv"
)

// StoredHandler 处理存储型XSS请求
func StoredHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewPageData2(7, 11, "")
		html := ""
		db := database.DB

		// 处理删除请求
		if r.Method == http.MethodGet && r.URL.Query().Get("id") != "" {
			id := r.URL.Query().Get("id")
			idNum, err := strconv.Atoi(id)
			if err == nil {
				_, err := db.Exec("DELETE FROM message WHERE id=$1", idNum)
				if err != nil {
					html += "<p id='op_notice'>删除失败,请重试并检查数据库是否还好!</p>"
				} else {
					// 重定向
					http.Redirect(w, r, "/vul/xss/xss_stored", http.StatusFound)
					return
				}
			}
		}

		// 处理POST请求，添加新留言
		if r.Method == http.MethodPost {
			message := r.FormValue("message")
			if message != "" {
				// 修正PostgreSQL占位符
				_, err := db.Exec("INSERT INTO message(content, time) VALUES($1, now())", message)
				if err != nil {
					html += "<p>数据库出现异常，提交失败！" + err.Error() + "</p>"
				}
			}
		}

		// 获取所有留言，按时间降序排列
		rows, err := db.Query("SELECT id, content, time FROM message ORDER BY time DESC")
		if err != nil {
			html += "<p>获取留言列表失败！" + err.Error() + "</p>"
		} else {
			defer rows.Close()
			for rows.Next() {
				var id int
				var content string
				var time string
				err = rows.Scan(&id, &content, &time)
				if err != nil {
					continue
				}
				// XSS漏洞点：直接输出用户输入的内容
				html += fmt.Sprintf("<p class='con'>%s</p><p class='time'>%s</p><a href='/vul/xss/xss_stored?id=%d'>删除</a><br><br>", content, time, id)
			}
		}

		data.HtmlMsg = template.HTML(html)
		renderer.RenderPage(w, "xss/xss_stored.html", data)
	}
}
