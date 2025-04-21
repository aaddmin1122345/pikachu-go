package sqli

import (
	"database/sql"
	"fmt"
	"net/http"
	"pikachu-go/templates"

	_ "github.com/lib/pq"
)

// SqliBlindBHandler 布尔型盲注
func SqliBlindBHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从 URL 获取输入
		input := r.URL.Query().Get("input")
		result := "请输入一个值进行注入测试。"

		// 使用输入值执行 SQL 查询，模拟 SQL 注入
		if input != "" {
			// 模拟 SQL 注入的查询语句
			db, err := sql.Open("postgres", "user=pgsql password=pgsql dbname=pikachu-go sslmode=disable")
			if err != nil {
				http.Error(w, "数据库连接失败", http.StatusInternalServerError)
				return
			}
			defer db.Close()

			// 执行注入查询
			query := fmt.Sprintf("SELECT 1 FROM users WHERE username = '%s' AND password = '123456'", input)
			rows, err := db.Query(query)
			if err != nil {
				http.Error(w, "查询错误", http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			// 根据查询结果判断
			if rows.Next() {
				result = "True, 用户名存在！"
			} else {
				result = "False, 用户名不存在。"
			}
		}

		data := templates.NewPageData2(3, 5, result)
		renderer.RenderPage(w, "sqli/sqli_blind_b.html", data)
	}
}
