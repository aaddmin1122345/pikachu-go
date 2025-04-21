package rce

import (
	"html/template"
	"net/http"
	"os/exec"
	"pikachu-go/templates"
)

func RceEvalHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		output := ""
		if r.Method == http.MethodPost {
			cmd := r.FormValue("cmd")
			// 安全起见，我们只允许简单命令
			out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
			if err != nil {
				output = err.Error() + "\n"
			}
			output += string(out)
		}
		data := templates.NewPageData2(50, 53, template.HTMLEscapeString(output))
		renderer.RenderPage(w, "rce/rce_eval.html", data)
	}
}
