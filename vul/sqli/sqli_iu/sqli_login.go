package sqliiu

import (
	"database/sql"
	"fmt"
	"net/http"
	"pikachu-go/templates"

	_ "github.com/lib/pq"
)

// SqliLoginHandler 登录页面
func SqliLoginHandler(renderer templates.Renderer) http.HandlerFunc {
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

			query := fmt.Sprintf("SELECT id FROM users WHERE username = '%s' AND password = '%s'", username, password)
			row := db.QueryRow(query)

			var id int
			if err := row.Scan(&id); err == nil {
				http.SetCookie(w, &http.Cookie{
					Name:  "login_user",
					Value: username,
					Path:  "/",
				})
				http.Redirect(w, r, "/vul/sqli/sqli_mem", http.StatusFound)
				return
			}
			msg = "登录失败"
		}

		data := templates.NewPageData2(3, 12, msg)
		renderer.RenderPage(w, "sqli/sqli_login.html", data)
	}
}
