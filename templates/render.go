package templates

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type PageData struct {
	PikaRoot string
	Active   []string
	HtmlMsg  template.HTML
	Extra    map[string]interface{}
}

type Renderer interface {
	RenderPage(w http.ResponseWriter, contentTemplate string, data PageData) error
}

type TemplateRenderer struct {
	templates *template.Template
}

func NewTemplateRenderer() (Renderer, error) {
	// tmpl := template.New("")
	tmpl := template.New("").Funcs(template.FuncMap{
		"list": func(args ...string) []string {
			return args
		},
		"split": func(s, sep string) []string {
			return strings.Split(s, sep)
		},
	})

	err := filepath.Walk("templates", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".html") {
			return nil
		}
		relPath := strings.TrimPrefix(path, "templates/")
		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		_, parseErr := tmpl.New(relPath).Parse(string(content))
		return parseErr
	})

	if err != nil {
		return nil, err
	}

	for _, t := range tmpl.Templates() {
		log.Println("加载模板名：", t.Name())
	}

	return &TemplateRenderer{templates: tmpl}, nil
}

func (tr *TemplateRenderer) RenderPage(w http.ResponseWriter, contentTemplate string, data PageData) error {
	log.Println("开始渲染页面:", contentTemplate)

	if err := tr.templates.ExecuteTemplate(w, "header.html", data); err != nil {
		log.Println("执行 header 模板错误：", err)
		return err
	}

	if err := tr.templates.ExecuteTemplate(w, contentTemplate, data); err != nil {
		log.Println("执行内容模板错误：", err)
		return err
	}

	if err := tr.templates.ExecuteTemplate(w, "footer.html", data); err != nil {
		log.Println("执行 footer 模板错误：", err)
		return err
	}

	// log.Println("成功渲染页面:", contentTemplate)
	return nil
}

// 旧版已弃用，请使用 NewPageData2 替代
// func NewPageData(...) 已不推荐使用

func NewPageData2(activeMain, activeSub int, htmlMsg string) PageData {
	active := make([]string, 130)
	if activeMain >= 0 && activeMain < len(active) {
		active[activeMain] = "active open"
	}
	if activeSub >= 0 && activeSub < len(active) {
		active[activeSub] = "active"
	}
	return PageData{
		PikaRoot: "/",
		Active:   active,
		HtmlMsg:  template.HTML(htmlMsg),
		Extra:    make(map[string]interface{}),
	}
}
