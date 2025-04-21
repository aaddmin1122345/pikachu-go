package sqli

import (
	"net/http"
	"pikachu-go/templates"
)

// SqliWidebyteHandler 字节宽度注入
func SqliWidebyteHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 示例：模拟字节宽度注入
		result := "字节宽度注入结果："
		input := r.URL.Query().Get("input")

		if input != "" {
			result = "注入的字节：" + input
		}

		data := templates.NewPageData2(3, 11, result)
		renderer.RenderPage(w, "sqli/sqli_widebyte.html", data)
	}
}
