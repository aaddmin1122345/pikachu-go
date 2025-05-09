package sqli

import (
	"fmt"
	"net/http"
	"pikachu-go/database"
	"pikachu-go/templates"
)

// SqliIDHandler 数字型注入漏洞演示
func SqliIDHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		htmlMsg := ""

		// 获取用户输入
		id := r.URL.Query().Get("id")

		if id != "" {
			// 使用全局数据库连接
			db := database.DB
			if db == nil {
				htmlMsg = "数据库连接失败"
			} else {
				// 故意保留SQL注入漏洞 - 直接拼接用户输入到SQL语句中
				query := fmt.Sprintf("SELECT username, email FROM users WHERE id = %s", id)

				rows, err := db.Query(query)
				if err != nil {
					htmlMsg = fmt.Sprintf("<p class='notice'>查询错误: %s</p>", err.Error())
				} else {
					defer rows.Close()

					found := false
					for rows.Next() {
						found = true
						var username, email string
						if err := rows.Scan(&username, &email); err != nil {
							htmlMsg = fmt.Sprintf("<p class='notice'>读取数据失败: %s</p>", err.Error())
						} else {
							htmlMsg = fmt.Sprintf("<p class='notice'>Hello, %s <br />Your email is: %s</p>",
								username, email)
						}
					}

					if !found {
						htmlMsg = "<p class='notice'>您输入的user id不存在，请重新输入！</p>"
					}
				}
			}
		}

		data := templates.NewPageData2(35, 37, htmlMsg)
		renderer.RenderPage(w, "sqli/sqli_id.html", data)
	}
}
