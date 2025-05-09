package rce

import (
	"html/template"
	"net/http"
	"os/exec"
	"pikachu-go/templates"
	"runtime"
)

// RcePingHandler 实现命令注入漏洞演示
func RcePingHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result := ""

		if r.Method == http.MethodPost {
			ip := r.FormValue("ipaddress")
			if ip != "" {
				// 模拟PHP版本中的命令注入漏洞
				// 故意直接拼接用户输入到命令中
				var cmd *exec.Cmd

				if runtime.GOOS == "windows" {
					// Windows系统
					cmd = exec.Command("cmd", "/c", "ping "+ip)
				} else {
					// Linux/Unix系统
					cmd = exec.Command("sh", "-c", "ping -c 4 "+ip)
				}

				// 执行命令并获取输出
				output, err := cmd.CombinedOutput()
				if err != nil {
					result = err.Error() + "\n"
				}
				result += string(output)
			}
		}

		data := templates.PageData{
			Active:  make([]string, 130),
			HtmlMsg: template.HTML("<pre>" + template.HTMLEscapeString(result) + "</pre>"),
		}
		data.Active[50] = "active open"
		data.Active[52] = "active"

		renderer.RenderPage(w, "rce/rce_ping.html", data)
	}
}
