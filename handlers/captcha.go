package handlers

import (
	"net/http"
	"pikachu-go/utils"
)

// CaptchaHandler 生成验证码图片并存储在会话中
func CaptchaHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 生成验证码并获取验证码字符串
		captchaText := utils.GenerateCaptcha(w)

		// 将验证码存储在会话中
		utils.GlobalSessions.SetSessionData(w, r, "vcode", captchaText)

		// 同时也设置一个Cookie，用于演示验证码绕过漏洞
		http.SetCookie(w, &http.Cookie{
			Name:  "bf[vcode]",
			Value: captchaText,
			Path:  "/",
		})
	}
}
