package unserilization

import (
	"html/template"
	"net/http"
	"pikachu-go/templates"
)

// UnserHandler 接收用户输入并显示（模拟 PHP unserialize 入口）
func UnserHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := ""

		if r.Method == http.MethodPost {
			raw := r.FormValue("txt")
			if raw == "" {
				msg = `<p style="color:red;">请输入序列化字符串</p>`
			} else {
				msg = `<p class="notice">收到模拟的序列化字符串（未执行反序列化，仅展示）：</p><pre>` +
					template.HTMLEscapeString(raw) + "</pre>"
			}
		}

		data := templates.NewPageData2(64, 66, msg)
		renderer.RenderPage(w, "unserilization/unser.html", data)
	}
}
