package sqliheader

import (
	"fmt"
	"net/http"
	"pikachu-go/database"
	"pikachu-go/templates"
)

// SqliHeaderLoginHandler 模拟基于 header 的注入登录
func SqliHeaderLoginHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从 Header 获取 username
		username := r.Header.Get("username")
		result := "请设置 HTTP Header 中的 username 来进行登录测试。"

		if username != "" {
			// 使用全局数据库连接
			db := database.DB
			if db == nil {
				http.Error(w, "数据库连接失败", http.StatusInternalServerError)
				return
			}

			// 故意保留SQL注入漏洞，直接拼接SQL
			query := fmt.Sprintf("SELECT id, username, password FROM users WHERE username = '%s'", username)
			row := db.QueryRow(query)

			var id int
			var uname, pwd string
			if err := row.Scan(&id, &uname, &pwd); err == nil {
				result = fmt.Sprintf("欢迎你，%s！", uname)
			} else {
				result = "登录失败，用户不存在。"
			}
		}

		data := templates.NewPageData2(35, 43, result)
		renderer.RenderPage(w, "sqli/sqli_header/sqli_header_login.html", data)
	}
}
