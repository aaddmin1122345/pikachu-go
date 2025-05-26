package sqli

import (
	"fmt"
	"net/http"
	"pikachu-go/database"
	"pikachu-go/templates"
)

// SqliBlindBHandler 布尔型盲注
func SqliBlindBHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从 URL 获取输入
		input := r.URL.Query().Get("input")
		result := "请输入一个值进行注入测试。"

		// 使用输入值执行 SQL 查询，模拟 SQL 注入
		if input != "" {
			// 使用全局数据库连接
			db := database.DB
			if db == nil {
				http.Error(w, "数据库连接失败", http.StatusInternalServerError)
				return
			}

			// 故意保留SQL注入漏洞，直接拼接SQL
			query := fmt.Sprintf("SELECT 1 FROM member WHERE username = '%s' AND password = '123456'", input)
			rows, err := db.Query(query)
			if err != nil {
				result = `<p style="color:red">查询错误: ` + err.Error() + `</p>`
			} else {
				defer rows.Close()

				// 根据查询结果判断
				if rows.Next() {
					result = "True, 用户名存在！"
				} else {
					result = "False, 用户名不存在。"
				}
			}
		}

		data := templates.NewPageData2(35, 44, result)
		renderer.RenderPage(w, "sqli/sqli_blind_b.html", data)
	}
}
