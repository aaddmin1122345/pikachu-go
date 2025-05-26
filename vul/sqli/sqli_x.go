package sqli

import (
	"fmt"
	"net/http"
	"pikachu-go/database"
	"pikachu-go/templates"
)

// SqliXHandler 其他 SQL 注入类型
func SqliXHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取用户输入
		input := r.URL.Query().Get("input")
		result := "请输入一个值进行测试。"

		if input != "" {
			// 使用全局数据库连接
			db := database.DB
			if db == nil {
				http.Error(w, "数据库连接失败", http.StatusInternalServerError)
				return
			}

			// 故意保留SQL注入漏洞，直接拼接SQL
			query := fmt.Sprintf("SELECT id, username, password FROM member WHERE username='%s'", input)
			rows, err := db.Query(query)
			if err != nil {
				http.Error(w, "查询错误", http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			result = "查询结果："
			for rows.Next() {
				var id int
				var username string
				var password string
				if err := rows.Scan(&id, &username, &password); err != nil {
					http.Error(w, "数据读取错误", http.StatusInternalServerError)
					return
				}
				result += fmt.Sprintf("ID: %d, 用户名: %s, 密码: %s<br>", id, username, password)
			}
		}

		data := templates.NewPageData2(35, 40, result)
		renderer.RenderPage(w, "sqli/sqli_x.html", data)
	}
}
