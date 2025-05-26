package templateinj

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os/exec"
	"pikachu-go/templates"
)

// TemplateInjHandler 渲染模板注入漏洞概述页面
func TemplateInjHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewPageData2(116, 117, "")
		renderer.RenderPage(w, "templateinj/templateinj.html", data)
	}
}

// TemplateInjTestHandler 处理模板注入测试请求
func TemplateInjTestHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		htmlMsg := ""
		var result string

		// 模拟用户输入的个性签名
		note := ""
		if r.Method == "POST" {
			note = r.FormValue("tpl") // 沿用tpl参数名，但含义变为note
		}

		// 业务变量
		user := "pikachu"
		date := "2025-05-26 12:00:00"

		// 如果用户没有输入，使用默认模板
		tplToParse := note
		if tplToParse == "" && r.Method != "POST" {
			// Using backticks for cleaner string literal
			tplToParse = `欢迎{{.user}}，现在是{{.date}}`
		}

		// exec函数 - **危险操作，仅用于靶场演示**
		exec := func(cmd string) string {
			out, err := exec.Command("sh", "-c", cmd).CombinedOutput()
			if err != nil {
				return fmt.Sprintf("命令执行错误: %s\n输出: %s", err.Error(), string(out))
			}
			return string(out)
		}

		// 使用text/template演示漏洞
		t, err := template.New("test").Funcs(template.FuncMap{"exec": exec}).Parse(tplToParse)
		if err != nil {
			result = "模板解析错误: " + err.Error()
		} else {
			data := map[string]interface{}{
				"user": user,
				"date": date,
			}
			var buf bytes.Buffer
			err = t.Execute(&buf, data)
			if err != nil {
				result = "模板执行错误: " + err.Error()
			} else {
				result = buf.String()
			}
		}

		if result != "" {
			htmlMsg = "<div class='alert alert-success'><pre>" + result + "</pre></div>"
		}

		data := templates.NewPageData2(116, 118, htmlMsg)
		// 将用户输入的note也传回前端，方便回显
		data.Extra["tpl"] = note
		data.Extra["user"] = user // 传递user和date到Extra，以便前端示例显示
		data.Extra["date"] = date

		// 注册exec函数到主模板的Extra中，使其在渲染主模板时可用
		data.Extra["exec"] = exec

		renderer.RenderPage(w, "templateinj/templateinj_test.html", data)
	}
}
