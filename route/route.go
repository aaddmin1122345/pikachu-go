package route

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"pikachu-go/database"
	"pikachu-go/templates"
	"pikachu-go/vul/apiunauth"
	"pikachu-go/vul/burteforce"
	"pikachu-go/vul/dir"
	"pikachu-go/vul/fileinclude"
	"pikachu-go/vul/infoleak"
	"pikachu-go/vul/rce"
	"pikachu-go/vul/sqli"
	sqliiu "pikachu-go/vul/sqli/sqli_iu"
	"pikachu-go/vul/ssrf"
	"pikachu-go/vul/ssti"
	"pikachu-go/vul/unsafedownload"
	"pikachu-go/vul/unsafeupload"
	"pikachu-go/vul/urlredirect"
	"pikachu-go/vul/xss"
	"pikachu-go/vul/xss/xssblind"
	"pikachu-go/vul/xss/xsspost"
)

func indexHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		active := make([]string, 130)
		active[0] = "active open"

		// 使用数据库连接
		db := database.DB
		htmlMsg := ""

		if db == nil {
			log.Println("数据库打开失败")
			htmlMsg = `<p><a href="/install" style="color:red;">提示：无法读取数据库！</a></p>`
		} else {
			if err := db.Ping(); err != nil {
				log.Println("数据库连接失败：", err)
				htmlMsg = `<p><a href="/install" style="color:red;">提示:无法连接数据库!</a></p>`
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
		htmlMsg := "<p>欢迎使用 Pikachu-Go 安装向导！</p>"
		active := make([]string, 130)
		active[0] = "active open"
		data := templates.PageData{
			Active:  active,
			HtmlMsg: template.HTML(htmlMsg),
		}
		// 检查 install.flag 文件是否存在
		_, err := os.Stat("./install.flag")
		installFileExists := os.IsNotExist(err)
		if installFileExists {
			// 调用数据库初始化
			database.InitializeDatabase()

			// 更新现有用户的密码为MD5格式
			database.UpdateExistingPasswords()

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

func InitRoutes() {
	// // 初始化数据库
	database.InitializeDatabase()

	// // 更新现有用户的密码为MD5格式
	database.UpdateExistingPasswords()

	renderer, err := templates.NewTemplateRenderer()
	if err != nil {
		log.Fatal("加载模板失败：", err)
	}

	// 首页
	http.HandleFunc("/", indexHandler(renderer))

	// // 静态资源处理（放在路由处理之前以确保优先级）
	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	// vulFs := http.FileServer(http.Dir("./vul"))
	// http.Handle("/vul/", http.StripPrefix("/vul/", vulFs))

	// ========== URL 重定向模块 ==========
	http.HandleFunc("/vul/urlredirect/urlredirect", urlredirect.URLRedirectHandler(renderer))
	http.HandleFunc("/vul/urlredirect/unsafere", urlredirect.UnsafeReHandler(renderer))

	// ========== RCE 远程命令执行模块 ==========
	http.HandleFunc("/vul/rce/rce", rce.RceIndexHandler(renderer))
	// http.HandleFunc("/vul/rce/rce_eval", rce.RceEvalHandler(renderer))
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
	http.HandleFunc("/vul/xss/xsspost/post_login", xsspost.PostLoginHandler(renderer))
	http.HandleFunc("/vul/xss/xssblind/xss_blind", xssblind.XssBlindHandler(renderer))
	http.HandleFunc("/vul/xss/xssblind/admin", xssblind.AdminHandler(renderer))
	http.HandleFunc("/vul/xss/xssblind/admin_login", xssblind.AdminLoginHandler(renderer))

	http.HandleFunc("/vul/dir/dir", dir.DirHandler(renderer))
	http.HandleFunc("/vul/dir/dir_list", dir.DirListHandler(renderer))

	// ======== 信息泄露模块 ========
	http.HandleFunc("/vul/infoleak/infoleak", infoleak.InfoleakHandler(renderer))
	http.HandleFunc("/vul/infoleak/findabc", infoleak.FindABCHandler(renderer))
	http.HandleFunc("/vul/infoleak/abc", infoleak.ABCHandler(renderer))

	// 暂时下线，后面在优化
	// // ======== 越权模块 ========
	// http.HandleFunc("/vul/overpermission/op", overpermission.OpHandler(renderer))
	// http.HandleFunc("/vul/overpermission/op1/op1_login", op1.Op1LoginHandler(renderer))
	// http.HandleFunc("/vul/overpermission/op1/op1_mem", op1.Op1MemHandler(renderer))
	// http.HandleFunc("/vul/overpermission/op2/op2_login", op2.Op2LoginHandler(renderer))
	// http.HandleFunc("/vul/overpermission/op2/op2_user", op2.Op2UserHandler(renderer))
	// http.HandleFunc("/vul/overpermission/op2/op2_admin", op2.Op2AdminHandler(renderer))
	// http.HandleFunc("/vul/overpermission/op2/op2_admin_edit", op2.Op2AdminEditHandler(renderer))

	http.HandleFunc("/vul/unsafedownload/unsafedownload", unsafedownload.UnsafedownloadHandler(renderer))
	http.HandleFunc("/vul/unsafedownload/down_nba", unsafedownload.DownNbaHandler(renderer))
	http.HandleFunc("/vul/unsafedownload/execdownload", unsafedownload.ExecDownloadHandler())
	http.HandleFunc("/vul/unsafeupload/upload", unsafeupload.UploadHandler(renderer))
	http.HandleFunc("/vul/unsafeupload/clientcheck", unsafeupload.ClientcheckHandler(renderer))
	http.HandleFunc("/vul/unsafeupload/servercheck", unsafeupload.ServercheckHandler(renderer))
	http.HandleFunc("/vul/unsafeupload/getimagesize", unsafeupload.GetimagesizeHandler(renderer))

	// Go自带的xml标注库不支持xml注入，只能使用第三方库，暂时下线，后续改为html模板注入漏洞
	// http.HandleFunc("/vul/xxe/xxe", xxe.XxeHandler(renderer))
	// http.HandleFunc("/vul/xxe/xxe_1", xxe.Xxe1Handler(renderer))

	http.HandleFunc("/vul/fileinclude/fileinclude", fileinclude.FileIncludeHandler(renderer))
	http.HandleFunc("/vul/fileinclude/fi_local", fileinclude.FiLocalHandler(renderer))
	http.HandleFunc("/vul/fileinclude/fi_remote", fileinclude.FiRemoteHandler(renderer))

	http.HandleFunc("/vul/ssrf/ssrf", ssrf.SsrfHandler(renderer))
	http.HandleFunc("/vul/ssrf/ssrf_curl", ssrf.SsrfCurlHandler(renderer))
	// http.HandleFunc("/vul/ssrf/ssrf_fgc", ssrf.SsrfFgcHandler(renderer))

	http.HandleFunc("/vul/ssrf/ssrf_info/info1", ssrf.Info1Handler())
	http.HandleFunc("/vul/ssrf/ssrf_info/info2", ssrf.Info2Handler())

	// 经测试，目前无法实现这个漏洞，所以暂时下线
	// http.HandleFunc("/vul/unserilization/unserilization", unserilization.UnserilizationHandler(renderer))
	// http.HandleFunc("/vul/unserilization/unser", unserilization.UnserHandler(renderer))

	// ========== 暴力破解模块 ==========
	http.HandleFunc("/vul/burteforce/burteforce", burteforce.BurteforceHandler(renderer))
	http.HandleFunc("/vul/burteforce/bf_form", burteforce.BfFormHandler(renderer))
	http.HandleFunc("/vul/burteforce/bf_server", burteforce.BfServerHandler(renderer))
	http.HandleFunc("/vul/burteforce/bf_client", burteforce.BfClientHandler(renderer))
	http.HandleFunc("/vul/burteforce/vcode", burteforce.VcodeHandler)
	http.HandleFunc("/vul/burteforce/bf_token", burteforce.BfTokenHandler(renderer))

	// csrf的全部下线，还待实现
	// http.HandleFunc("/vul/csrf/csrf", csrf.CsrfHandler(renderer))
	// http.HandleFunc("/vul/csrf/login", csrf.CsrfLoginHandler(renderer))
	// http.HandleFunc("/vul/csrf/debug_login", csrf.DebugLoginHandler(renderer))

	// http.HandleFunc("/vul/csrf/csrfget/csrf_get_login", csrfget.CsrfGetLoginHandler(renderer))
	// http.HandleFunc("/vul/csrf/csrfget/csrf_get_edit", csrfget.CsrfGetEditHandler(renderer))
	// http.HandleFunc("/vul/csrf/csrfget/csrf_get", csrfget.CsrfGetHandler(renderer))
	// // CSRF GET攻击演示页面
	// http.HandleFunc("/vul/csrf/csrfget/hack.html", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "pikachu-go/vul/csrf/csrfget/hack.html")
	// })

	// http.HandleFunc("/vul/csrf/csrfpost/csrf_post_login", csrfpost.CsrfPostLoginHandler(renderer))
	// http.HandleFunc("/vul/csrf/csrfpost/csrf_post_edit", csrfpost.CsrfPostEditHandler(renderer))
	// http.HandleFunc("/vul/csrf/csrfpost/csrf_post", csrfpost.CsrfPostHandler(renderer))
	// // CSRF POST攻击演示页面
	// http.HandleFunc("/vul/csrf/csrfpost/hack.html", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "pikachu-go/vul/csrf/csrfpost/hack.html")
	// })

	// http.HandleFunc("/vul/csrf/csrftoken/token_get_login", csrftoken.TokenGetLoginHandler(renderer))
	// http.HandleFunc("/vul/csrf/csrftoken/token_get_edit", csrftoken.TokenGetEditHandler(renderer))
	// http.HandleFunc("/vul/csrf/csrftoken/token_get", csrftoken.TokenGetHandler(renderer))
	// // CSRF Token攻击演示页面
	// http.HandleFunc("/vul/csrf/csrftoken/hack.html", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "pikachu-go/vul/csrf/csrftoken/hack.html")
	// })

	// 安装路由
	http.HandleFunc("/install", installHandler(renderer))

	// sql注入路由
	http.HandleFunc("/vul/sqli/sqli", sqli.SqliHandler(renderer))
	http.HandleFunc("/vul/sqli/sqli_blind_b", sqli.SqliBlindBHandler(renderer))
	http.HandleFunc("/vul/sqli/sqli_blind_t", sqli.SqliBlindTHandler(renderer))
	http.HandleFunc("/vul/sqli/sqli_del", sqli.SqliDelHandler(renderer))
	http.HandleFunc("/vul/sqli/sqli_id", sqli.SqliIDHandler(renderer))
	http.HandleFunc("/vul/sqli/sqli_search", sqli.SqliSearchHandler(renderer))
	http.HandleFunc("/vul/sqli/sqli_str", sqli.SqliStrHandler(renderer))
	// 暂时下线，有问题
	// http.HandleFunc("/vul/sqli/sqli_widebyte", sqli.SqliWidebyteHandler(renderer))
	http.HandleFunc("/vul/sqli/sqli_x", sqli.SqliXHandler(renderer))
	// 暂时下线了这两个，暂时有问题
	// http.HandleFunc("/vul/sqli/sqli_header/sqli_header", sqliheader.SqliHeaderHandler(renderer))
	// http.HandleFunc("/vul/sqli/sqli_header/sqli_header_login", sqliheader.SqliHeaderLoginHandler(renderer))
	http.HandleFunc("/vul/sqli/sqli_iu/sqli_login", sqliiu.SqliLoginHandler(renderer))
	http.HandleFunc("/vul/sqli/sqli_iu/sqli_mem", sqliiu.SqliMemHandler(renderer))
	http.HandleFunc("/vul/sqli/sqli_iu/sqli_reg", sqliiu.SqliRegHandler(renderer))
	http.HandleFunc("/vul/sqli/sqli_iu/sqli_edit", sqliiu.SqliEditHandler(renderer))

	// // 初始化数据库
	// database.InitializeDatabase()

	// // 更新现有用户的密码为MD5格式
	// database.UpdateExistingPasswords()

	// ========== API未授权访问模块 ==========
	http.HandleFunc("/vul/apiunauth/apiunauth", apiunauth.ApiUnauthHandler(renderer))
	http.HandleFunc("/vul/apiunauth/apiunauth_demo", apiunauth.ApiUnauthDemoHandler(renderer))
	http.HandleFunc("/vul/apiunauth/api/users", apiunauth.ApiUsersHandler())
	http.HandleFunc("/vul/apiunauth/api/delete", apiunauth.ApiDeleteUserHandler())

	// ========== 模板注入模块 ==========
	http.HandleFunc("/vul/ssti/ssti", ssti.TemplateInjHandler(renderer))
	http.HandleFunc("/vul/ssti/exec", ssti.TemplateInjTestHandler(renderer))
}
