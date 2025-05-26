package sqli

import (
	"fmt"
	"net/http"
	"pikachu-go/database"
	"pikachu-go/templates"
	"time"
)

// SqliBlindTHandler 时间型盲注
func SqliBlindTHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取用户输入
		input := r.URL.Query().Get("input")
		result := "请输入一个值进行测试。"

		// 如果用户输入不为空，则模拟时间型盲注
		if input != "" {
			// 使用全局数据库连接
			db := database.DB
			if db == nil {
				http.Error(w, "数据库连接失败", http.StatusInternalServerError)
				return
			}

			// 构造时间型盲注查询，模拟长时间查询
			query := fmt.Sprintf("SELECT 1 FROM member WHERE username = '%s' AND password = '123456'", input)

			// 开始时间记录
			start := time.Now()

			// 执行查询
			_, err := db.Query(query)
			if err != nil {
				http.Error(w, "查询错误", http.StatusInternalServerError)
				return
			}

			// 计算查询所用时间
			duration := time.Since(start)

			// 根据查询时间判断是否存在用户名
			if duration.Seconds() > 2 {
				result = "True, 延迟符合时间要求，用户名存在。"
			} else {
				result = "False, 延迟不符合要求，用户名不存在。"
			}
		}

		data := templates.NewPageData2(35, 45, result)
		renderer.RenderPage(w, "sqli/sqli_blind_t.html", data)
	}
}
