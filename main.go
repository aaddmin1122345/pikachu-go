package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// PageData 定义页面渲染所需要的数据
type PageData struct {
	PikaRoot string
	Active   []string
	HtmlMsg  template.HTML
}

// loadTemplates 从 templates 目录下加载模板
func loadTemplates() (*template.Template, error) {
	// 模板按照 header, index, footer 拼接，模板中可以互相引用数据
	tmpl, err := template.ParseFiles(
		filepath.Join("templates", "header.html"),
		filepath.Join("templates", "index.html"),
		filepath.Join("templates", "footer.html"),
	)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

// renderPage 组合模板并输出完整页面
func renderPage(w http.ResponseWriter, data PageData, tmpl *template.Template) {
	// 模板执行的顺序为 header -> index -> footer
	err := tmpl.ExecuteTemplate(w, "header.html", data)
	if err != nil {
		log.Println("执行 header 模板错误：", err)
		return
	}
	err = tmpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		log.Println("执行 index 模板错误：", err)
		return
	}
	err = tmpl.ExecuteTemplate(w, "footer.html", data)
	if err != nil {
		log.Println("执行 footer 模板错误：", err)
		return
	}
}

// indexHandler 处理首页请求
func indexHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 定义 PikaRoot，默认空字符串，可根据需要修改
		pikaRoot := ""

		// 定义 Active 数组，按照 PHP 中 index.php 的设置，将第 0 项设为 "active open"，其他项为空
		active := make([]string, 130)
		active[0] = "active open"
		// 其它项可按需要设置

		// 尝试连接 SQLite 数据库文件 "pikachu.db"
		db, err := sql.Open("sqlite3", "./pikachu.db")
		htmlMsg := ""
		if err != nil {
			// 数据库连接失败，提示安装
			htmlMsg = `<p>
        <a href="install.php" style="color:red;">
        提示:欢迎使用, pikachu还没有初始化，点击进行初始化安装!
        </a>
    </p>`
		} else {
			// 测试数据库连接是否正常
			err = db.Ping()
			if err != nil {
				htmlMsg = `<p>
        <a href="install.php" style="color:red;">
        提示:欢迎使用, pikachu还没有初始化，点击进行初始化安装!
        </a>
    </p>`
			}
		}
		// 关闭数据库连接
		if db != nil {
			db.Close()
		}

		data := PageData{
			PikaRoot: pikaRoot,
			Active:   active,
			HtmlMsg:  template.HTML(htmlMsg),
		}
		renderPage(w, data, tmpl)
	}
}

func main() {
	tmpl, err := loadTemplates()
	if err != nil {
		log.Fatal("加载模板错误：", err)
	}
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	http.Handle("/vul/", http.StripPrefix("/vul/", http.FileServer(http.Dir("./vul/"))))

	http.HandleFunc("/", indexHandler(tmpl))
	// 监听 8080 端口
	log.Println("服务器启动在 http://localhost:8888")
	err = http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
