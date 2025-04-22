package sqliiu

import (
	"database/sql"
	"fmt"
	"net/http"
	"pikachu-go/templates"

	_ "github.com/lib/pq"
)

// SqliRegHandler 注册页面（含注入）
func SqliRegHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := ""
		if r.Method == http.MethodPost {
			username := r.FormValue("username")
			password := r.FormValue("password")

			db, err := sql.Open("postgres", "user=pgsql password=pgsql dbname=pikachu-go sslmode=disable")
			if err != nil {
				http.Error(w, "数据库连接失败", http.StatusInternalServerError)
				return
			}
			defer db.Close()

			query := fmt.Sprintf("INSERT INTO users(username, password) VALUES('%s', '%s')", username, password)
			_, err = db.Exec(query)
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
