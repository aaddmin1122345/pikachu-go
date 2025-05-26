package apiunauth

import (
	"fmt"
	"net/http"
	"pikachu-go/database"
	"pikachu-go/templates"
	"strconv"
)

// Member 表示用户数据结构
type Member struct {
	ID       int
	Username string
	Password string
	Sex      string
	Phonenum string
	Address  string
	Email    string
}

// authMiddleware 检查认证的中间件
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 检查Authorization头
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// 这里简单检查token是否存在，实际应用中应该验证token的有效性
		if authHeader != "Bearer valid-token" {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

// ApiUnauthHandler 处理API未授权访问漏洞的概述页面
func ApiUnauthHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewPageData2(111, 112, "")
		renderer.RenderPage(w, "apiunauth/apiunauth.html", data)
	}
}

// ApiUnauthDemoHandler 处理API未授权访问漏洞的演示页面
func ApiUnauthDemoHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewPageData2(111, 113, "")
		renderer.RenderPage(w, "apiunauth/apiunauth_demo.html", data)
	}
}

// ApiUsersHandler 处理用户API请求 - GET请求不需要鉴权，其他请求需要鉴权
func ApiUsersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取用户列表或单个用户
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Missing id parameter", http.StatusBadRequest)
			return
		}

		// 获取单个用户
		userID, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		var member Member
		err = database.DB.QueryRow("SELECT * FROM member WHERE id = $1", userID).Scan(
			&member.ID, &member.Username, &member.Password, &member.Sex, &member.Phonenum, &member.Address, &member.Email)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		response := fmt.Sprintf("ID: %d\nUsername: %s\nPassword: %s\nSex: %s\nPhone: %s\nAddress: %s\nEmail: %s",
			member.ID, member.Username, member.Password, member.Sex, member.Phonenum, member.Address, member.Email)
		w.Write([]byte(response))

	}
}

// ApiDeleteUserHandler 处理删除用户的API请求
func ApiDeleteUserHandler() http.HandlerFunc {
	return authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		id := r.FormValue("id")
		if id == "" {
			http.Error(w, "Missing id parameter", http.StatusBadRequest)
			return
		}

		userID, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		_, err = database.DB.Exec("DELETE FROM member WHERE id = $1", userID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		w.Write([]byte("User deleted successfully"))
	})
}
