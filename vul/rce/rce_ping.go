package rce

import (
	"html/template"
	"net/http"
	"os/exec"
	"pikachu-go/templates"
	"regexp"
)

// 这个代码还没写好，
func RcePingHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result := ""

		if r.Method == http.MethodPost {
			ip := r.FormValue("ip")

			// 简单校验只允许 IP/域名格式
			validInput := regexp.MustCompile(`^[a-zA-Z0-9\.\-]+$`).MatchString
			if validInput(ip) {
				cmd := exec.Command("ping", "-c", "1", ip)
				out, err := cmd.CombinedOutput()
				if err != nil {
					result += err.Error() + "\n"
				}
				result += string(out)
			} else {
				result = "输入非法，疑似注入命令，已拦截。\n"
			}
		}

		data := templates.NewPageData2(50, 52, template.HTMLEscapeString(result))
		renderer.RenderPage(w, "rce/rce_ping.html", data)
	}
}
