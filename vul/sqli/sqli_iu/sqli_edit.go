package sqliiu

import (
	"database/sql"
	"fmt"
	"net/http"
	"pikachu-go/templates"

	_ "github.com/mattn/go-sqlite3"
)

// SqliEditHandler 用户编辑操作 SQL 注入
func SqliEditHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取输入字段（例如昵称）
		nickname := r.URL.Query().Get("nickname")
		result := "请输入要修改的昵称。"

		if nickname != "" {
			// 执行 SQL 注入进行修改
			db, err := sql.Open("sqlite3", "./pikachu.db")
			if err != nil {
				http.Error(w, "数据库连接失败", http.StatusInternalServerError)
				return
			}
			defer db.Close()

			// 模拟 SQL 注入进行修改操作
			query := fmt.Sprintf("UPDATE users SET username='%s' WHERE id=1", nickname)
			_, err = db.Exec(query)
			if err != nil {
				http.Error(w, "数据库更新失败", http.StatusInternalServerError)
				return
			}

			result = fmt.Sprintf("成功修改昵称为：%s", nickname)
		}

		data := templates.NewPageData2(35, 41, result)
		renderer.RenderPage(w, "sqli/sqli_iu/sqli_edit.html", data)
	}
}
