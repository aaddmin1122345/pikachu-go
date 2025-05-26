package sqliheader

import (
	"fmt"
	"net/http"
	"pikachu-go/database"
	"pikachu-go/templates"
)

// SqliHeaderHandler 请求头部 SQL 注入
func SqliHeaderHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取自定义请求头
		headerValue := r.Header.Get("X-User")
		result := "请在请求头中添加 X-User 字段进行测试。"

		if headerValue != "" {
			// 使用全局数据库连接
			db := database.DB
			if db == nil {
				http.Error(w, "数据库连接失败", http.StatusInternalServerError)
				return
			}

			// 故意保留SQL注入漏洞，直接拼接SQL
			query := fmt.Sprintf("SELECT id, username, password FROM users WHERE username='%s'", headerValue)
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

		data := templates.NewPageData2(35, 43, result)
		renderer.RenderPage(w, "sqli/sqli_header/sqli_header.html", data)
	}
}
