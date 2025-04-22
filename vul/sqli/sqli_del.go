package sqli

import (
	"database/sql"
	"fmt"
	"net/http"
	"pikachu-go/templates"

	_ "github.com/lib/pq"
)

// SqliDelHandler SQL 删除操作注入
func SqliDelHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取要删除的 ID
		id := r.URL.Query().Get("id")
		result := "请输入删除的 ID。"

		// 模拟通过 SQL 注入删除操作
		if id != "" {
			db, err := sql.Open("postgres", "user=pgsql password=pgsql dbname=pikachu-go sslmode=disable")
			if err != nil {
				http.Error(w, "数据库连接失败", http.StatusInternalServerError)
				return
			}
			defer db.Close()

			// 使用 SQL 注入删除用户
			query := fmt.Sprintf("DELETE FROM users WHERE id = '%s'", id)
			_, err = db.Exec(query)
			if err != nil {
				http.Error(w, "删除失败", http.StatusInternalServerError)
				return
			}

			result = fmt.Sprintf("成功删除 ID：%s", id)
		}

		data := templates.NewPageData2(35, 42, result)
		renderer.RenderPage(w, "sqli/sqli_del.html", data)
	}
}
