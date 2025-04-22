package utils

import (
	"math/rand"
	"net/http"
	"sync"
)

// 简单会话管理器，用于存储验证码等会话数据
type SessionManager struct {
	sessions map[string]map[string]interface{}
	mutex    sync.RWMutex
}

var GlobalSessions = NewSessionManager()

// 创建新的会话管理器
func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]map[string]interface{}),
	}
}

// 从请求中获取会话ID（使用cookie）
func (m *SessionManager) getSessionID(r *http.Request) string {
	cookie, err := r.Cookie("PIKASESSION")
	if err != nil {
		return ""
	}
	return cookie.Value
}

// 设置会话ID（通过cookie）
func (m *SessionManager) setSessionID(w http.ResponseWriter, id string) {
	cookie := &http.Cookie{
		Name:     "PIKASESSION",
		Value:    id,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

// 获取会话，如果不存在则创建
func (m *SessionManager) GetSession(w http.ResponseWriter, r *http.Request) map[string]interface{} {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	id := m.getSessionID(r)
	if id == "" || m.sessions[id] == nil {
		// 生成新会话ID
		id = generateSessionID()
		m.sessions[id] = make(map[string]interface{})
		m.setSessionID(w, id)
	}

	return m.sessions[id]
}

// 设置会话数据
func (m *SessionManager) SetSessionData(w http.ResponseWriter, r *http.Request, key string, value interface{}) {
	session := m.GetSession(w, r)

	m.mutex.Lock()
	defer m.mutex.Unlock()

	session[key] = value
}

// 获取会话数据
func (m *SessionManager) GetSessionData(r *http.Request, key string) (interface{}, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	id := m.getSessionID(r)
	if id == "" || m.sessions[id] == nil {
		return nil, false
	}

	val, ok := m.sessions[id][key]
	return val, ok
}

// 生成唯一会话ID
func generateSessionID() string {
	// 简化实现，实际应使用更安全的随机生成方法
	return "session_" + randString(32)
}

// 生成随机字符串
func randString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
