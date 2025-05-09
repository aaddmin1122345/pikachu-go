package unserilization

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"pikachu-go/templates"
	"strings"
)

// PikachuObject 模拟可被反序列化的对象
type PikachuObject struct {
	Test string `json:"test"`
}

// UnserHandler 反序列化漏洞演示
func UnserHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		htmlMsg := ""

		if r.Method == http.MethodPost {
			serializedData := r.FormValue("o")
			if serializedData == "" {
				htmlMsg = "<p>大兄弟,来点劲爆点儿的!</p>"
			} else {
				// 尝试解析JSON数据（模拟反序列化）
				// 注意: 这里故意实现了一个有漏洞的反序列化过程
				var obj PikachuObject

				// 危险的方法: 直接将用户输入反序列化为对象
				err := json.Unmarshal([]byte(serializedData), &obj)

				if err != nil {
					htmlMsg = "<p>大兄弟,来点劲爆点儿的!</p>"
				} else {
					// 故意输出未过滤的内容，可能导致XSS
					// 这里模拟PHP中未过滤的输出，直接把用户可控制的数据输出到页面
					htmlMsg = fmt.Sprintf("<p>%s</p>", obj.Test)
				}

				// 提示: 尝试输入 {"test":"<script>alert('xss')</script>"}
				if !strings.Contains(serializedData, "xss") && !strings.Contains(htmlMsg, "xss") {
					htmlMsg += "<p><small>提示: 尝试输入包含XSS的JSON数据，例如 {\"test\":\"&lt;script&gt;alert('xss')&lt;/script&gt;\"}</small></p>"
				}
			}
		}

		data := templates.PageData{
			Active:  make([]string, 130),
			HtmlMsg: template.HTML(htmlMsg),
		}
		data.Active[90] = "active open"
		data.Active[92] = "active"

		renderer.RenderPage(w, "unserilization/unser.html", data)
	}
}
