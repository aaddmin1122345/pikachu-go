package xxe

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"pikachu-go/templates"
)

// 简单通用结构用于解析 xml 内容
type AnyXML struct {
	XMLName xml.Name
	Content string `xml:",innerxml"`
}

// Xxe1Handler 支持用户提交 XML 文本并解析（模拟 XXE 行为）
func Xxe1Handler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		htmlMsg := ""

		if r.Method == http.MethodPost {
			xmlText := r.FormValue("xmlcontent")
			if xmlText == "" {
				htmlMsg = `<p style="color:red;">请输入 XML 内容</p>`
			} else {
				var parsed AnyXML
				err := xml.Unmarshal([]byte(xmlText), &parsed)
				if err != nil {
					htmlMsg = fmt.Sprintf(`<p style="color:red;">解析失败：%s</p>`, err.Error())
				} else {
					htmlMsg = fmt.Sprintf(`<p style="color:green;">解析成功！XML 内容：</p><pre>%s</pre>`, xmlText)
				}
			}
		}

		data := templates.NewPageData2(95, 97, htmlMsg)
		renderer.RenderPage(w, "xxe/xxe_1.html", data)
	}
}
