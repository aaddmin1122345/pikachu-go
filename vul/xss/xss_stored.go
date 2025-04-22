package xss

import (
	"fmt"
	"net/http"
	"pikachu-go/database"
	"pikachu-go/templates"
	"strconv"
)

func StoredHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		html := ""
		db := database.DB

		// 处理删除请求
		if r.Method == http.MethodGet && r.URL.Query().Get("id") != "" {
			id := r.URL.Query().Get("id")
			idNum, err := strconv.Atoi(id)
			if err == nil {
				// SQL注入漏洞（故意保留，作为彩蛋）
				query := fmt.Sprintf("DELETE FROM message WHERE id=%d", idNum)
				_, err := db.Exec(query)
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
				_, err := db.Exec("INSERT INTO message(content,time) VALUES(?,now())", message)
				if err != nil {
					html += "<p>数据库出现异常，提交失败！</p>"
				}
			}
		}

		// 获取所有留言
		rows, err := db.Query("SELECT * FROM message")
		if err != nil {
			html += "<p>获取留言列表失败！</p>"
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
				// 改进留言显示的HTML格式
				html += fmt.Sprintf(`<div class="message-item">
                    <p class="con">%s</p>
                    <p class="message-meta">时间: %s | <a href='/vul/xss/xss_stored?id=%d'>删除</a></p>
                </div>`, content, time, id)
			}
		}

		data := templates.NewPageData2(7, 11, html)
		renderer.RenderPage(w, "xss/xss_stored.html", data)
	}
}
