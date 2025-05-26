package rce

import (
	"html/template"
	"net/http"
	"os/exec"
	"pikachu-go/templates"
	"runtime"
)

// RcePingHandler 处理ping命令注入漏洞展示
func RcePingHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result := ""

		if r.Method == http.MethodPost && r.FormValue("submit") != "" {
			ip := r.FormValue("ipaddress")
			if ip != "" {
				// 直接使用用户输入进行命令拼接，模拟命令注入漏洞
				if runtime.GOOS == "windows" {
					// Windows系统
					cmd := exec.Command("cmd", "/c", "ping "+ip)
					output, err := cmd.CombinedOutput()
					if err != nil {
						result = err.Error() + "\n"
					}
					result += string(output)
				} else {
					// Linux/Unix系统
					cmd := exec.Command("sh", "-c", "ping -c 4 "+ip)
					output, err := cmd.CombinedOutput()
					if err != nil {
						result = err.Error() + "\n"
					}
					result += string(output)
				}
			}
		}

		// 使用NewPageData2创建页面数据，使用正确的菜单ID
		htmlOutput := "<pre>" + template.HTMLEscapeString(result) + "</pre>"
		data := templates.NewPageData2(50, 52, htmlOutput)

		renderer.RenderPage(w, "rce/rce_ping.html", data)
	}
}
