package sqli

import (
	"database/sql"
	"fmt"
	"net/http"
	"pikachu-go/templates"

	_ "github.com/lib/pq"
)

// SqliIDHandler 获取 ID 注入
func SqliIDHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取用户输入
		id := r.URL.Query().Get("id")
		result := "未找到数据"

		// 模拟 SQL 注入获取 ID
		if id != "" {
			db, err := sql.Open("postgres", "user=pgsql password=pgsql dbname=pikachu-go sslmode=disable")
			if err != nil {
				http.Error(w, "数据库连接失败", http.StatusInternalServerError)
				return
			}
			defer db.Close()

			// 使用输入的 ID 执行查询
			query := fmt.Sprintf("SELECT username FROM users WHERE id = '%s'", id)
			rows, err := db.Query(query)
			if err != nil {
				http.Error(w, "查询错误", http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			// 获取查询结果
			if rows.Next() {
				var username string
				err := rows.Scan(&username)
				if err != nil {
					http.Error(w, "读取数据失败", http.StatusInternalServerError)
					return
				}
				result = fmt.Sprintf("查询到的用户名：%s", username)
			}
		}

		data := templates.NewPageData2(35, 37, result)
		renderer.RenderPage(w, "sqli/sqli_id.html", data)
	}
}
