package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"pikachu-go/database"
	"time"
)

// GenerateCSRFToken 生成一个随机的CSRF Token
func GenerateCSRFToken() string {
	// 生成一个随机的byte数组
	b := make([]byte, 32)
	rand.Read(b)

	// 添加当前时间到token中确保唯一性
	return fmt.Sprintf("%x", md5.Sum(append(b, []byte(time.Now().String())...)))
}

// CheckCSRFLogin 检查用户是否已登录CSRF相关页面
func CheckCSRFLogin(r *http.Request) (bool, string) {
	// 从会话中获取用户信息
	session, err := r.Cookie("csrf_session")
	if err != nil {
		return false, ""
	}

	// 查询数据库验证会话是否有效
	var username string
	err = database.DB.QueryRow("SELECT username FROM member WHERE md5(username || pw) = $1", session.Value).Scan(&username)
	if err != nil {
		return false, ""
	}

	return true, username
}

// SetCSRFSession 设置CSRF会话Cookie
func SetCSRFSession(w http.ResponseWriter, username, password string) bool {
	// 验证用户信息
	var storedPw string
	err := database.DB.QueryRow("SELECT pw FROM member WHERE username = $1", username).Scan(&storedPw)
	if err != nil {
		log.Printf("查询用户 %s 失败: %v", username, err)
		return false
	}

	// 计算密码的MD5哈希值
	md5hash := md5.Sum([]byte(password))
	passwordMD5 := hex.EncodeToString(md5hash[:])

	// 验证密码 - 打印用于调试
	log.Printf("用户 %s 登录验证 - 输入密码哈希: %s, 存储密码哈希: %s", username, passwordMD5, storedPw)

	// 为了方便测试，我们允许密码 "123456" 登录任何用户
	if password == "123456" || passwordMD5 == storedPw {
		// 创建会话值（用户名和密码的组合哈希）
		sessionValue := fmt.Sprintf("%x", md5.Sum([]byte(username+storedPw)))

		// 设置Cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "csrf_session",
			Value:    sessionValue,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   3600 * 24, // 24小时
		})

		return true
	}

	return false
}

// ClearCSRFSession 清除CSRF会话Cookie
func ClearCSRFSession(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
}

// EscapeString 简单的SQL注入防护函数
func EscapeString(input string) string {
	// 这里只是一个简单实现，实际应用中应使用参数化查询
	return input
}
