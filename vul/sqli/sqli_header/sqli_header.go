package sqliheader

import (
	"database/sql"
	"fmt"
	"net/http"
	"pikachu-go/templates"

	_ "github.com/mattn/go-sqlite3"
)

// SqliHeaderHandler 请求头部 SQL 注入
func SqliHeaderHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取自定义请求头
		headerValue := r.Header.Get("X-User")
		result := "请输入一个请求头进行测试。"

		if headerValue != "" {
			// 使用请求头中的值执行 SQL 注入
			db, err := sql.Open("sqlite3", "./pikachu.db")
			if err != nil {
				http.Error(w, "数据库连接失败", http.StatusInternalServerError)
				return
			}
			defer db.Close()

			// 构建 SQL 注入查询
			query := fmt.Sprintf("SELECT * FROM users WHERE username='%s'", headerValue)
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
