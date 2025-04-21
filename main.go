package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"pikachu-go/templates"
	"pikachu-go/vul/burteforce"
	"pikachu-go/vul/dir"
	"pikachu-go/vul/fileinclude"
	"pikachu-go/vul/infoleak"
	"pikachu-go/vul/rce"
	"pikachu-go/vul/ssrf"
	"pikachu-go/vul/unsafedownload"
	"pikachu-go/vul/unsafeupload"
	"pikachu-go/vul/unserilization"
	"pikachu-go/vul/urlredirect"
	"pikachu-go/vul/xss"
	"pikachu-go/vul/xss/xssblind"
	"pikachu-go/vul/xss/xsspost"
	"pikachu-go/vul/xxe"
)

func indexHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		active := make([]string, 130)
		active[0] = "active open"

		db, err := sql.Open("sqlite3", "./pikachu.db")
		htmlMsg := ""

		if err != nil {
			log.Println("数据库打开失败：", err)
			htmlMsg = `<p><a href="install.php" style="color:red;">提示:欢迎使用, pikachu还没有初始化，点击进行初始化安装!</a></p>`
		} else {
			defer db.Close()
			if err = db.Ping(); err != nil {
				log.Println("数据库连接失败：", err)
				htmlMsg = `<p><a href="install.php" style="color:red;">提示:欢迎使用, pikachu还没有初始化，点击进行初始化安装!</a></p>`
			}
		}

		data := templates.PageData{
			Active:  active,
			HtmlMsg: template.HTML(htmlMsg),
		}

		if err := renderer.RenderPage(w, "index.html", data); err != nil {
			log.Println("渲染页面失败：", err)
			http.Error(w, "渲染页面失败："+err.Error(), http.StatusInternalServerError)
		}
	}
}

func main() {
	renderer, err := templates.NewTemplateRenderer()
	if err != nil {
		log.Fatal("加载模板失败：", err)
	}

	// 首页
	http.HandleFunc("/", indexHandler(renderer))

	// ========== URL 重定向模块 ==========
	http.HandleFunc("/vul/urlredirect/urlredirect", urlredirect.URLRedirectHandler(renderer))
	http.HandleFunc("/vul/urlredirect/unsafere", urlredirect.UnsafeReHandler(renderer))

	// ========== RCE 远程命令执行模块 ==========
	http.HandleFunc("/vul/rce/rce", rce.RceIndexHandler(renderer))
	http.HandleFunc("/vul/rce/rce_eval", rce.RceEvalHandler(renderer))
	http.HandleFunc("/vul/rce/rce_ping", rce.RcePingHandler(renderer))

	// ========== XSS 模块 ==========
	http.HandleFunc("/vul/xss/xss", xss.XssIndexHandler(renderer))
	http.HandleFunc("/vul/xss/xss_01", xss.RenderXssVariant(renderer, "xss_01"))
	http.HandleFunc("/vul/xss/xss_02", xss.RenderXssVariant(renderer, "xss_02"))
	http.HandleFunc("/vul/xss/xss_03", xss.RenderXssVariant(renderer, "xss_03"))
	http.HandleFunc("/vul/xss/xss_04", xss.RenderXssVariant(renderer, "xss_04"))
	http.HandleFunc("/vul/xss/xss_reflected_get", xss.ReflectedGetHandler(renderer))
	http.HandleFunc("/vul/xss/xss_stored", xss.StoredHandler(renderer))
	http.HandleFunc("/vul/xss/xss_dom", xss.DomXssHandler(renderer))
	http.HandleFunc("/vul/xss/xss_dom_x", xss.DomXssXHandler(renderer))
	http.HandleFunc("/vul/xss/fixxss", xss.FixXssHandler(renderer))
	http.HandleFunc("/vul/xss/xsspost/xss_reflected_post", xsspost.XssReflectedPostHandler(renderer))
	http.HandleFunc("/vul/xss/xssblind/xss_blind", xssblind.XssBlindHandler(renderer))
	http.HandleFunc("/vul/xss/xssblind/admin", xssblind.AdminHandler(renderer))
	http.HandleFunc("/vul/xss/xssblind/admin_login", xssblind.AdminLoginHandler(renderer))

	http.HandleFunc("/vul/dir/dir", dir.DirHandler(renderer))
	http.HandleFunc("/vul/dir/dir_list", dir.DirListHandler(renderer))

	// ======== 信息泄露模块 ========
	http.HandleFunc("/vul/infoleak/infoleak", infoleak.InfoleakHandler(renderer))
	http.HandleFunc("/vul/infoleak/findabc", infoleak.FindABCHandler(renderer))
	http.HandleFunc("/vul/infoleak/abc", infoleak.ABCHandler(renderer))

	// ========== 静态资源兜底：放最后 ==========
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	http.Handle("/vul/", http.StripPrefix("/vul/", http.FileServer(http.Dir("./vul/"))))

	http.HandleFunc("/vul/unsafedownload/unsafedownload", unsafedownload.UnsafedownloadHandler(renderer))
	http.HandleFunc("/vul/unsafedownload/down_nba", unsafedownload.DownNbaHandler(renderer))
	http.HandleFunc("/vul/unsafedownload/execdownload", unsafedownload.ExecDownloadHandler())

	http.HandleFunc("/vul/unsafeupload/upload", unsafeupload.UploadHandler(renderer))
	http.HandleFunc("/vul/unsafeupload/clientcheck", unsafeupload.ClientcheckHandler(renderer))
	http.HandleFunc("/vul/unsafeupload/servercheck", unsafeupload.ServercheckHandler(renderer))
	http.HandleFunc("/vul/unsafeupload/getimagesize", unsafeupload.GetimagesizeHandler(renderer))
	http.HandleFunc("/vul/xxe/xxe", xxe.XxeHandler(renderer))
	http.HandleFunc("/vul/xxe/xxe_1", xxe.Xxe1Handler(renderer))
	http.HandleFunc("/vul/fileinclude/fileinclude", fileinclude.FileIncludeHandler(renderer))
	http.HandleFunc("/vul/fileinclude/fi_local", fileinclude.FiLocalHandler(renderer))
	http.HandleFunc("/vul/fileinclude/fi_remote", fileinclude.FiRemoteHandler(renderer))
	http.HandleFunc("/vul/ssrf/ssrf", ssrf.SsrfHandler(renderer))
	http.HandleFunc("/vul/ssrf/ssrf_curl", ssrf.SsrfCurlHandler(renderer))
	http.HandleFunc("/vul/ssrf/ssrf_fgc", ssrf.SsrfFgcHandler(renderer))
	http.HandleFunc("/vul/unserilization/unserilization", unserilization.UnserilizationHandler(renderer))
	http.HandleFunc("/vul/unserilization/unser", unserilization.UnserHandler(renderer))
	http.HandleFunc("/vul/burteforce/burteforce", burteforce.BurteforceHandler(renderer))
	http.HandleFunc("/vul/burteforce/bf_form", burteforce.BfFormHandler(renderer))
	http.HandleFunc("/vul/burteforce/bf_server", burteforce.BfServerHandler(renderer))

	http.HandleFunc("/vul/burteforce/bf_client", burteforce.BfClientHandler(renderer))
	http.HandleFunc("/vul/burteforce/bf_token", burteforce.BfTokenHandler(renderer))

	// 启动服务
	log.Println("服务器启动于 http://localhost:8888")
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatal("服务器启动失败: ", err)
	}
}
