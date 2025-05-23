package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"pikachu-go/database"
	"pikachu-go/handlers"
	"pikachu-go/templates"
	"pikachu-go/vul/burteforce"
	"pikachu-go/vul/csrf"
	"pikachu-go/vul/csrf/csrfget"
	"pikachu-go/vul/csrf/csrfpost"
	"pikachu-go/vul/csrf/csrftoken"
	"pikachu-go/vul/dir"
	"pikachu-go/vul/fileinclude"
	"pikachu-go/vul/infoleak"
	"pikachu-go/vul/overpermission"
	"pikachu-go/vul/overpermission/op1"
	"pikachu-go/vul/overpermission/op2"
	"pikachu-go/vul/rce"
	"pikachu-go/vul/sqli"
	sqliheader "pikachu-go/vul/sqli/sqli_header"
	sqliiu "pikachu-go/vul/sqli/sqli_iu"
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

		// 使用数据库连接
		db := database.DB
		htmlMsg := ""

		if db == nil {
			log.Println("数据库打开失败")
			htmlMsg = `<p><a href="/install" style="color:red;">提示:欢迎使用, pikachu还没有初始化，点击进行初始化安装!</a></p>`
		} else {
			if err := db.Ping(); err != nil {
				log.Println("数据库连接失败：", err)
				htmlMsg = `<p><a href="/install" style="color:red;">提示:欢迎使用, pikachu还没有初始化，点击进行初始化安装!</a></p>`
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

func installHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		htmlMsg := "ok"
		active := make([]string, 130)
		active[0] = "active open"
		data := templates.PageData{
			Active:  active,
			HtmlMsg: template.HTML(htmlMsg),
		}
		// 检查 install.flag 文件是否存在
		_, err := os.Stat("./install.flag")
		installFileExists := os.IsNotExist(err)

		// 如果 install.flag 或 pikachu.db 文件不存在，则初始化数据库
		if installFileExists {
			// 调用数据库初始化
			database.InitializeDatabase()

			// 创建 install.flag 文件
			_, err := os.Create("./install.flag")
			if err != nil {
				log.Println("创建 install.flag 文件失败：", err)
				http.Error(w, "初始化失败，请重试", http.StatusInternalServerError)
				return
			}

			// 返回数据库初始化成功
			renderer.RenderPage(w, "install_success.html", data)
		} else {
			// 数据库已经初始化，返回提示
			renderer.RenderPage(w, "install_already.html", data)
		}
	}
}

func main() {
	// 初始化数据库
	database.InitializeDatabase()

	renderer, err := templates.NewTemplateRenderer()
	if err != nil {
		log.Fatal("加载模板失败：", err)
	}

	// 首页
	http.HandleFunc("/", indexHandler(renderer))

	// 验证码路由
	http.HandleFunc("/captcha", handlers.CaptchaHandler())

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

	// ======== 越权模块 ========
	http.HandleFunc("/vul/overpermission/op", overpermission.OpHandler(renderer))
	http.HandleFunc("/vul/overpermission/op1/op1_login", op1.Op1LoginHandler(renderer))
	http.HandleFunc("/vul/overpermission/op1/op1_mem", op1.Op1MemHandler(renderer))
	http.HandleFunc("/vul/overpermission/op2/op2_login", op2.Op2LoginHandler(renderer))
	http.HandleFunc("/vul/overpermission/op2/op2_user", op2.Op2UserHandler(renderer))
	http.HandleFunc("/vul/overpermission/op2/op2_admin", op2.Op2AdminHandler(renderer))
	http.HandleFunc("/vul/overpermission/op2/op2_admin_edit", op2.Op2AdminEditHandler(renderer))

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

	// ========== 暴力破解模块 ==========
	http.HandleFunc("/vul/burteforce/burteforce", burteforce.BurteforceHandler(renderer))
	http.HandleFunc("/vul/burteforce/bf_form", burteforce.BfFormHandler(renderer))
	http.HandleFunc("/vul/burteforce/bf_server", burteforce.BfServerHandler(renderer))
	http.HandleFunc("/vul/burteforce/bf_client", burteforce.BfClientHandler(renderer))
	http.HandleFunc("/vul/burteforce/vcode", burteforce.VcodeHandler)
	http.HandleFunc("/vul/burteforce/bf_token", burteforce.BfTokenHandler(renderer))

	http.HandleFunc("/vul/csrf/csrf", csrf.CsrfHandler(renderer))
	http.HandleFunc("/vul/csrf/csrfget/csrf_get_login", csrfget.CsrfGetLoginHandler(renderer))
	http.HandleFunc("/vul/csrf/csrfget/csrf_get_edit", csrfget.CsrfGetEditHandler())
	http.HandleFunc("/vul/csrf/csrfget/csrf_get", csrfget.CsrfGetHandler(renderer))
	http.HandleFunc("/vul/csrf/csrfpost/csrf_post_login", csrfpost.CsrfPostLoginHandler(renderer))
	http.HandleFunc("/vul/csrf/csrfpost/csrf_post_edit", csrfpost.CsrfPostEditHandler())
	http.HandleFunc("/vul/csrf/csrfpost/csrf_post", csrfpost.CsrfPostHandler(renderer))
	http.HandleFunc("/vul/csrf/csrftoken/token_get_login", csrftoken.TokenGetLoginHandler(renderer))
	http.HandleFunc("/vul/csrf/csrftoken/token_get_edit", csrftoken.TokenGetEditHandler())
	http.HandleFunc("/vul/csrf/csrftoken/token_get", csrftoken.TokenGetHandler(renderer))

	// 安装路由
	http.HandleFunc("/install", installHandler(renderer))

	http.HandleFunc("/vul/sqli/sqli", sqli.SqliHandler(renderer))
	http.HandleFunc("/vul/sqli/sqli_blind_b", sqli.SqliBlindBHandler(renderer))
	http.HandleFunc("/vul/sqli/sqli_blind_t", sqli.SqliBlindTHandler(renderer))
	http.HandleFunc("/vul/sqli/sqli_del", sqli.SqliDelHandler(renderer))
	http.HandleFunc("/vul/sqli/sqli_id", sqli.SqliIDHandler(renderer))
	http.HandleFunc("/vul/sqli/sqli_search", sqli.SqliSearchHandler(renderer))
	http.HandleFunc("/vul/sqli/sqli_str", sqli.SqliStrHandler(renderer))
	http.HandleFunc("/vul/sqli/sqli_widebyte", sqli.SqliWidebyteHandler(renderer))
	http.HandleFunc("/vul/sqli/sqli_x", sqli.SqliXHandler(renderer))
	http.HandleFunc("/vul/sqli/sqli_header/sqli_header", sqliheader.SqliHeaderHandler(renderer))
	http.HandleFunc("/vul/sqli/sqli_header/sqli_header_login", sqliheader.SqliHeaderLoginHandler(renderer))
	http.HandleFunc("/vul/sqli/sqli_iu/sqli_login", sqliiu.SqliLoginHandler(renderer))
	http.HandleFunc("/vul/sqli/sqli_iu/sqli_mem", sqliiu.SqliMemHandler(renderer))
	http.HandleFunc("/vul/sqli/sqli_iu/sqli_reg", sqliiu.SqliRegHandler(renderer))
	http.HandleFunc("/vul/sqli/sqli_iu/sqli_edit", sqliiu.SqliEditHandler(renderer))

	// 启动服务
	log.Println("服务器启动于 http://localhost:8888")
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatal("服务器启动失败: ", err)
	}
}
