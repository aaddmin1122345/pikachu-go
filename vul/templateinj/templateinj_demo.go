package templateinj

import (
	"bytes"
	"html/template"
	"net/http"
	"pikachu-go/templates"
)

// TemplateInjDemoHandler 处理模板注入漏洞的演示页面
func TemplateInjDemoHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		htmlMsg := ""
		var result string

		if r.Method == "POST" {
			tpl := r.FormValue("tpl")
			if tpl != "" {
				// 故意使用不安全的text/template演示漏洞
				t, err := template.New("test").Parse(tpl)
				if err != nil {
					result = "模板解析错误: " + err.Error()
				} else {
					data := map[string]interface{}{
						"User":   "pikachu",
						"Secret": "this_is_a_secret_key",
					}
					var buf bytes.Buffer
					err = t.Execute(&buf, data)
					if err != nil {
						result = "模板执行错误: " + err.Error()
					} else {
						result = buf.String()
					}
				}
			}
		}

		if result != "" {
			htmlMsg = "<div class='alert alert-success'><pre>" + result + "</pre></div>"
		}

		data := templates.NewPageData2(117, 118, htmlMsg)
		renderer.RenderPage(w, "templateinj/templateinj_demo.html", data)
	}
}
