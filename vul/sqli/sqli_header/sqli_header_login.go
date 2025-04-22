package sqliheader

import (
	"database/sql"
	"fmt"
	"net/http"
	"pikachu-go/templates"

	_ "github.com/lib/pq"
)

// SqliHeaderLoginHandler 模拟基于 header 的注入登录
func SqliHeaderLoginHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从 Header 获取 username
		username := r.Header.Get("username")
		result := "请设置 HTTP Header 中的 username 来进行登录测试。"

		if username != "" {
			db, err := sql.Open("postgres", "user=pgsql password=pgsql dbname=pikachu-go sslmode=disable")
			if err != nil {
				http.Error(w, "数据库连接失败", http.StatusInternalServerError)
				return
			}
			defer db.Close()

			// 模拟 SQL 注入
			query := fmt.Sprintf("SELECT * FROM users WHERE username = '%s'", username)
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
