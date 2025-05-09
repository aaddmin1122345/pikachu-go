package xxe

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"pikachu-go/templates"
	"strings"
)

// SimpleXML 用于解析XML的简单结构
type SimpleXML struct {
	XMLName xml.Name `xml:"root"`
	Content string   `xml:",innerxml"`
}

// Xxe1Handler 实现XXE漏洞演示
func Xxe1Handler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		htmlMsg := ""

		if r.Method == http.MethodPost {
			xmlContent := r.FormValue("xmlcontent")
			if xmlContent == "" {
				htmlMsg = `<p style="color:red;">请输入 XML 内容</p>`
			} else {
				// 使用临时文件实现XXE漏洞
				// 创建临时文件存放XML内容
				tmpFile, err := ioutil.TempFile("", "xxe_*.xml")
				if err != nil {
					htmlMsg = fmt.Sprintf(`<p style="color:red;">创建临时文件失败：%s</p>`, err.Error())
				} else {
					defer os.Remove(tmpFile.Name())

					// 写入XML内容
					if _, err := tmpFile.WriteString(xmlContent); err != nil {
						htmlMsg = fmt.Sprintf(`<p style="color:red;">写入XML失败：%s</p>`, err.Error())
					} else {
						tmpFile.Close()

						// 使用外部命令解析XML (这里模拟XXE漏洞)
						// 在实际系统中，这会导致服务器文件被读取
						var output []byte
						if strings.Contains(xmlContent, "SYSTEM") || strings.Contains(xmlContent, "ENTITY") {
							// 如果包含外部实体声明，尝试使用危险的方式解析
							// 注意：这只是模拟效果，实际上文件不会被读取
							if strings.Contains(xmlContent, "file:///etc/passwd") {
								// 模拟读取到/etc/passwd文件的内容
								output = []byte("模拟读取系统文件: /etc/passwd 的内容\nroot:x:0:0:root:/root:/bin/bash\ndaemon:x:1:1:daemon:/usr/sbin:/usr/sbin/nologin\n...")
							} else if match := strings.Contains(xmlContent, "file://"); match {
								// 模拟尝试读取其他文件
								output = []byte("检测到对系统文件的访问尝试，XXE漏洞利用成功!")
							} else {
								output = []byte("外部实体已声明，但未指定有效的文件路径")
							}
						} else {
							// 安全解析XML
							var parsed SimpleXML
							if err := xml.Unmarshal([]byte(xmlContent), &parsed); err != nil {
								output = []byte(fmt.Sprintf("XML解析失败: %s", err.Error()))
							} else {
								output = []byte(fmt.Sprintf("XML解析成功，内容: %s", parsed.Content))
							}
						}

						htmlMsg = fmt.Sprintf(`<p style="color:green;">解析完成！</p>
							<p>XML内容：</p><pre>%s</pre>
							<p>解析结果：</p><pre>%s</pre>`,
							template.HTMLEscapeString(xmlContent),
							template.HTMLEscapeString(string(output)))
					}
				}

				// 添加提示信息
				if !strings.Contains(xmlContent, "SYSTEM") && !strings.Contains(xmlContent, "ENTITY") {
					htmlMsg += `<p style="color:blue;">提示：尝试使用包含外部实体的XML，例如：
<pre>
&lt;?xml version="1.0" encoding="UTF-8"?&gt;
&lt;!DOCTYPE foo [ &lt;!ENTITY xxe SYSTEM "file:///etc/passwd"&gt; ]&gt;
&lt;root&gt;&lt;name&gt;&amp;xxe;&lt;/name&gt;&lt;/root&gt;
</pre></p>`
				}
			}
		}

		data := templates.PageData{
			Active:  make([]string, 130),
			HtmlMsg: template.HTML(htmlMsg),
		}
		data.Active[95] = "active open"
		data.Active[97] = "active"

		renderer.RenderPage(w, "xxe/xxe_1.html", data)
	}
}
