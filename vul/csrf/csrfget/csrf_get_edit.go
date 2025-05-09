package csrfget

import (
	"log"
	"net/http"
	"strings"
)

// CsrfGetEditHandler 通过GET请求修改昵称 - 故意存在CSRF漏洞
func CsrfGetEditHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取昵称参数
		nickname := r.URL.Query().Get("nickname")
		referrer := r.Referer()

		// 记录修改情况，便于排查
		if nickname != "" {
			log.Printf("CSRF GET演示: 昵称被修改为 %s, 来源: %s", nickname, referrer)

			// 如果昵称包含脚本标签，记录可能的XSS尝试
			if strings.Contains(strings.ToLower(nickname), "<script") {
				log.Printf("警告: 可能的XSS尝试被检测到 - %s", nickname)
			}
		}

		// 设置cookie - 故意不做任何CSRF防护
		http.SetCookie(w, &http.Cookie{
			Name:     "csrf_get_nick",
			Value:    nickname,
			Path:     "/",
			HttpOnly: false, // 允许JavaScript访问，增加风险
			MaxAge:   3600,  // 1小时过期
		})

		// 重定向回查看页面
		http.Redirect(w, r, "/vul/csrf/csrfget/csrf_get", http.StatusFound)
	}
}
