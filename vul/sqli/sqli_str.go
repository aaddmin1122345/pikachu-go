package sqli

import (
	"fmt"
	"net/http"
	"pikachu-go/database"
	"pikachu-go/templates"

	_ "github.com/lib/pq"
)

// SqliStrHandler 字符串注入
func SqliStrHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取用户输入
		str := r.URL.Query().Get("str")
		result := "请输入一个值进行测试。"

		// 执行字符串类型的注入
		if str != "" {
			// 使用全局数据库连接
			db := database.DB
			if db == nil {
				http.Error(w, "数据库连接失败", http.StatusInternalServerError)
				return
			}

			// 故意保留SQL注入漏洞，直接拼接SQL
			query := fmt.Sprintf("SELECT id, username, email FROM member WHERE username = '%s'", str)
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
				var id int
				var username string
				var email string
				err := rows.Scan(&id, &username, &email)
				if err != nil {
					http.Error(w, "读取数据失败", http.StatusInternalServerError)
					return
				}
				result += fmt.Sprintf("ID: %d, 用户名: %s, 邮箱: %s<br>", id, username, email)
			}

			if !found {
				result += "未找到匹配结果"
			}
		}

		data := templates.NewPageData2(35, 38, result)
		renderer.RenderPage(w, "sqli/sqli_str.html", data)
	}
}
