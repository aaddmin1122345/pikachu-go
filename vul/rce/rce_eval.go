package rce

import (
	"net/http"
	"os/exec"
	"pikachu-go/templates"
)

// 定义允许执行的命令白名单
// var allowedCommands = map[string]bool{
// 	"ls":       true,
// 	"dir":      true,
// 	"pwd":      true,
// 	"whoami":   true,
// 	"date":     true,
// 	"uname":    true,
// 	"echo":     true,
// 	"cat":      true,
// 	"type":     true,
// 	"hostname": true,
// }

// // 用于判断命令是否在白名单中
// func isCommandAllowed(cmd string) bool {
// 	// 提取命令的第一个部分(即主命令)
// 	cmdParts := strings.Fields(cmd)
// 	if len(cmdParts) == 0 {
// 		return false
// 	}

// 	return allowedCommands[cmdParts[0]]
// }

// // 验证命令是否安全(不包含危险字符)
// func isCommandSafe(cmd string) bool {
// 	// 禁止特殊字符和命令链接符号
// 	dangerousPatterns := []string{
// 		";", "|", "&", "&&", "||", ">", ">>", "<", "$(", "`",
// 		"\\", "\n", "\r", "eval", "exec", "system",
// 	}

// 	for _, pattern := range dangerousPatterns {
// 		if strings.Contains(cmd, pattern) {
// 			return false
// 		}
// 	}

// 	// 进一步使用正则表达式检查安全性
// 	safePattern := regexp.MustCompile(`^[a-zA-Z0-9\s._-]+$`)
// 	return safePattern.MatchString(cmd)
// }

// RceEvalHandler 处理eval命令执行漏洞展示
func RceEvalHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		html := ""

		if r.Method == http.MethodPost && r.FormValue("submit") != "" {
			txt := r.FormValue("txt")
			if txt != "" {
				// 在Go中模拟PHP的eval执行
				// 这里我们简单地尝试执行用户输入作为shell命令
				// 无论成功与否都返回相同的信息，以模拟原PHP行为
				cmd := exec.Command("sh", "-c", txt)
				cmd.CombinedOutput() // 忽略输出结果

				// 模拟PHP版本eval的回显信息
				html = "<p>你喜欢的字符还挺奇怪的!</p>"
			}
		}

		// 使用NewPageData2创建页面数据，使用正确的菜单ID
		data := templates.NewPageData2(50, 53, html)

		renderer.RenderPage(w, "rce/rce_eval.html", data)
	}
}
