package sqliiu

import (
	"fmt"
	"net/http"
	"pikachu-go/database"
	"pikachu-go/templates"
)

// SqliRegHandler 注册页面（含注入）
func SqliRegHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := ""
		if r.Method == http.MethodPost {
			username := r.FormValue("username")
			password := r.FormValue("password")

			// 使用全局数据库连接
			db := database.DB
			if db == nil {
				http.Error(w, "数据库连接失败", http.StatusInternalServerError)
				return
			}

			// 故意保留SQL注入漏洞，直接拼接SQL
			query := fmt.Sprintf("INSERT INTO member(username, password) VALUES('%s', '%s')", username, password)
			_, err := db.Exec(query)
			if err != nil {
				msg = "注册失败：" + err.Error()
			} else {
				msg = "注册成功"
			}
		}

		data := templates.NewPageData2(35, 41, msg)
		renderer.RenderPage(w, "sqli/sqli_iu/sqli_reg.html", data)
	}
}
