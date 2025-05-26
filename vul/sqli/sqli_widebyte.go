package sqli

import (
	"fmt"
	"net/http"
	"pikachu-go/database"
	"pikachu-go/templates"
)

// SqliWidebyteHandler 字节宽度注入
func SqliWidebyteHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result := "请输入一个值进行测试。"
		input := r.URL.Query().Get("input")

		if input != "" {
			// 使用全局数据库连接
			db := database.DB
			if db == nil {
				http.Error(w, "数据库连接失败", http.StatusInternalServerError)
				return
			}

			// 故意保留SQL注入漏洞，模拟宽字节注入场景
			// 在PostgreSQL中，我们假设这是一个容易受到宽字节注入的查询
			query := fmt.Sprintf("SELECT username FROM member WHERE username = '%s'", input)
			rows, err := db.Query(query)
			if err != nil {
				result = fmt.Sprintf("查询错误: %s", err.Error())
			} else {
				defer rows.Close()
				result = "查询结果: "

				for rows.Next() {
					var username string
					if err := rows.Scan(&username); err != nil {
						result += "数据读取错误"
					} else {
						result += username + " "
					}
				}
			}
		}

		data := templates.NewPageData2(35, 46, result)
		renderer.RenderPage(w, "sqli/sqli_widebyte.html", data)
	}
}
