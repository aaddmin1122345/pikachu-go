package sqli

import (
	"database/sql"
	"fmt"
	"net/http"
	"pikachu-go/templates"

	_ "github.com/lib/pq"
)

// SqliSearchHandler 查询注入
func SqliSearchHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取用户输入
		searchTerm := r.URL.Query().Get("search")
		result := "查询结果："

		// 执行查询操作
		if searchTerm != "" {
			db, err := sql.Open("postgres", "user=pgsql password=pgsql dbname=pikachu-go sslmode=disable")
			if err != nil {
				http.Error(w, "数据库连接失败", http.StatusInternalServerError)
				return
			}
			defer db.Close()

			// 模拟 SQL 查询注入
			query := fmt.Sprintf("SELECT username FROM users WHERE username LIKE '%%%s%%'", searchTerm)
			rows, err := db.Query(query)
			if err != nil {
				http.Error(w, "查询错误", http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			// 获取查询结果
			for rows.Next() {
				var username string
				err := rows.Scan(&username)
				if err != nil {
					http.Error(w, "读取数据失败", http.StatusInternalServerError)
					return
				}
				result += fmt.Sprintf("%s<br>", username)
			}
		}

		data := templates.NewPageData2(3, 9, result)
		renderer.RenderPage(w, "sqli/sqli_search.html", data)
	}
}
