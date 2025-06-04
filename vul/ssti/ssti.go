package ssti

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
		renderer.RenderPage(w, "ssti/ssti.html", data)
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
			note = r.FormValue("tpl")
		}

		// 业务变量
		user := "pikachu"
		date := "2025-05-26 12:00:00"

		// 如果用户没有输入，使用默认模板
		if note == "" && r.Method != "POST" {
			note = `欢迎{{.user}}，现在是{{.date}}`
		}

		// exec函数 模拟实际开发环境中用到了该命令
		exec := func(cmd string) string {
			out, err := exec.Command("sh", "-c", cmd).CombinedOutput()
			if err != nil {
				return fmt.Sprintf("命令执行错误: %s\n输出: %s", err.Error(), string(out))
			}
			return string(out)
		}

		// 使用text/template演示漏洞
		t, err := template.New("test").Funcs(template.FuncMap{"exec": exec}).Parse(note)
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
		// 传递给模板渲染
		data.Extra["tpl"] = note
		data.Extra["user"] = user
		data.Extra["date"] = date
		data.Extra["exec"] = exec

		renderer.RenderPage(w, "ssti/ssti_exec.html", data)
	}
}
