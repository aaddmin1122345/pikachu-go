package handlers

import (
	"log"
	"net/http"
	"pikachu-go/utils"
	"strings"
)

// CaptchaHandler 生成验证码图片并存储在会话中
func CaptchaHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 设置必要的头信息，确保图片不缓存
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		// 生成验证码并获取验证码字符串
		captchaText := utils.GenerateCaptcha(w)

		// 将验证码转换为小写并存储在会话中
		captchaText = strings.ToLower(captchaText)
		utils.GlobalSessions.SetSessionData(w, r, "vcode", captchaText)

		// 同时也设置一个Cookie，用于演示验证码绕过漏洞
		// 这是故意为之，模拟原PHP版本的行为
		http.SetCookie(w, &http.Cookie{
			Name:  "bf[vcode]",
			Value: captchaText,
			Path:  "/",
		})

		log.Printf("生成验证码: %s", captchaText)
	}
}
