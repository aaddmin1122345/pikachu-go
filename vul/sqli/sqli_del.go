package sqli

import (
	"fmt"
	"net/http"
	"pikachu-go/database"
	"pikachu-go/templates"
)

// SqliDelHandler SQL 删除操作注入
func SqliDelHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取要删除的 ID
		id := r.URL.Query().Get("id")
		result := "请输入要删除的 ID。"

		// 如果用户输入不为空，则模拟SQL注入
		if id != "" {
			// 使用全局数据库连接
			db := database.DB
			if db == nil {
				http.Error(w, "数据库连接失败", http.StatusInternalServerError)
				return
			}

			// 故意保留SQL注入漏洞，直接拼接SQL
			query := fmt.Sprintf("DELETE FROM member WHERE id = '%s'", id)
			_, err := db.Exec(query)
			if err != nil {
				http.Error(w, "删除失败: "+err.Error(), http.StatusInternalServerError)
				return
			}

			result = fmt.Sprintf("成功删除 ID：%s", id)
		}

		data := templates.NewPageData2(35, 42, result)
		renderer.RenderPage(w, "sqli/sqli_del.html", data)
	}
}
