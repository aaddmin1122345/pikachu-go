package sqliiu

import (
	"fmt"
	"net/http"
	"pikachu-go/database"
	"pikachu-go/templates"
)

// SqliEditHandler 用户编辑操作 SQL 注入
func SqliEditHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取输入字段（例如昵称）
		nickname := r.URL.Query().Get("nickname")
		result := "请输入要修改的昵称。"

		if nickname != "" {
			// 使用全局数据库连接
			db := database.DB
			if db == nil {
				http.Error(w, "数据库连接失败", http.StatusInternalServerError)
				return
			}

			// 故意保留SQL注入漏洞，直接拼接SQL
			query := fmt.Sprintf("UPDATE member SET username='%s' WHERE id=1", nickname)
			_, err := db.Exec(query)
			if err != nil {
				http.Error(w, "数据库更新失败: "+err.Error(), http.StatusInternalServerError)
				return
			}

			result = fmt.Sprintf("成功修改昵称为：%s", nickname)
		}

		data := templates.NewPageData2(35, 41, result)
		renderer.RenderPage(w, "sqli/sqli_iu/sqli_edit.html", data)
	}
}
