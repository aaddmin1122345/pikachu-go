package sqli

import (
	"fmt"
	"net/http"
	"pikachu-go/database"
	"pikachu-go/templates"
)

// SqliSearchHandler 查询注入
func SqliSearchHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取用户输入
		searchTerm := r.URL.Query().Get("search")
		result := "请输入查询条件。"

		// 执行查询操作
		if searchTerm != "" {
			// 使用全局数据库连接
			db := database.DB
			if db == nil {
				http.Error(w, "数据库连接失败", http.StatusInternalServerError)
				return
			}

			// 故意保留SQL注入漏洞，直接拼接SQL
			query := fmt.Sprintf("SELECT username FROM member WHERE username LIKE '%%%s%%'", searchTerm)
			rows, err := db.Query(query)
			if err != nil {
				http.Error(w, "查询错误: "+err.Error(), http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			// 获取查询结果
			found := false
			result = "查询结果: "
			for rows.Next() {
				found = true
				var username string
				err := rows.Scan(&username)
				if err != nil {
					http.Error(w, "读取数据失败", http.StatusInternalServerError)
					return
				}
				result += fmt.Sprintf("%s<br>", username)
			}

			if !found {
				result += "未找到匹配结果"
			}
		}

		data := templates.NewPageData2(35, 39, result)
		renderer.RenderPage(w, "sqli/sqli_search.html", data)
	}
}
