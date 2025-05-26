package xssblind

import (
	"fmt"
	"html/template"
	"net/http"
	"pikachu-go/database"
	"pikachu-go/templates"
	"time"
)

func XssBlindHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewPageData2(7, 14, "")
		htmlMsg := ""

		// 处理POST请求
		if r.Method == http.MethodPost {
			content := r.FormValue("content")
			name := r.FormValue("name")

			if content != "" {
				// 这里不过滤XSS，直接存入数据库，便于演示XSS盲打
				currentTime := time.Now().Format("2006-01-02 15:04:05")

				// 获取数据库连接
				db := database.DB
				if db == nil {
					htmlMsg = "<p>数据库连接错误</p>"
				} else {
					// 插入数据
					query := "INSERT INTO xssblind(time, content, name) VALUES($1, $2, $3)"
					_, err := db.Exec(query, currentTime, content, name)

					if err != nil {
						htmlMsg = fmt.Sprintf("<p>提交出现异常，请重新提交: %s</p>", err.Error())
					} else {
						htmlMsg = "<p>谢谢参与，阁下的看法我们已经收到!</p>"
					}
				}
			}
		}

		data.HtmlMsg = template.HTML(htmlMsg)
		renderer.RenderPage(w, "xss/xssblind/xss_blind.html", data)
	}
}
