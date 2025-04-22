package sqli

import (
	"database/sql"
	"fmt"
	"net/http"
	"pikachu-go/templates"

	_ "github.com/lib/pq"
)

// SqliStrHandler 字符串注入
func SqliStrHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取用户输入
		str := r.URL.Query().Get("str")
		result := "字符串注入结果："

		// 执行字符串类型的注入
		if str != "" {
			db, err := sql.Open("postgres", "user=pgsql password=pgsql dbname=pikachu-go sslmode=disable")
			if err != nil {
				http.Error(w, "数据库连接失败", http.StatusInternalServerError)
				return
			}
			defer db.Close()

			// 构造 SQL 查询
			query := fmt.Sprintf("SELECT * FROM users WHERE username = '%s'", str)
			rows, err := db.Query(query)
			if err != nil {
				http.Error(w, "查询错误", http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			// 获取查询结果
			for rows.Next() {
				var id int
				var username string
				var password string
				err := rows.Scan(&id, &username, &password)
				if err != nil {
					http.Error(w, "读取数据失败", http.StatusInternalServerError)
					return
				}
				result += fmt.Sprintf("ID: %d, 用户名: %s, 密码: %s<br>", id, username, password)
			}
		}

		data := templates.NewPageData2(35, 38, result)
		renderer.RenderPage(w, "sqli/sqli_str.html", data)
	}
}
