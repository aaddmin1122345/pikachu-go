package utils

import (
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// SessionManager 简单的会话管理器
type SessionManager struct {
	sessions map[string]map[string]interface{}
	mutex    sync.RWMutex
}

var GlobalSessions = NewSessionManager()

// 初始化随机数生成器
func init() {
	rand.Seed(time.Now().UnixNano())
}

// NewSessionManager 创建新的会话管理器
func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]map[string]interface{}),
	}
}

// GetSession 获取会话，如果不存在则创建
func (m *SessionManager) GetSession(w http.ResponseWriter, r *http.Request) map[string]interface{} {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	id := m.getSessionID(r)
	if id == "" || m.sessions[id] == nil {
		id = generateSessionID()
		m.sessions[id] = make(map[string]interface{})
		m.setSessionID(w, id)
	}

	return m.sessions[id]
}

// SetSessionData 设置会话数据
func (m *SessionManager) SetSessionData(w http.ResponseWriter, r *http.Request, key string, value interface{}) {
	session := m.GetSession(w, r)
	m.mutex.Lock()
	defer m.mutex.Unlock()
	session[key] = value
}

// GetSessionData 获取会话数据
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

// DeleteSessionData 删除会话数据
func (m *SessionManager) DeleteSessionData(r *http.Request, key string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	id := m.getSessionID(r)
	if id == "" || m.sessions[id] == nil {
		return
	}

	delete(m.sessions[id], key)
}

// 内部辅助函数
func (m *SessionManager) getSessionID(r *http.Request) string {
	cookie, err := r.Cookie("pikachu-session")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func (m *SessionManager) setSessionID(w http.ResponseWriter, id string) {
	cookie := &http.Cookie{
		Name:     "pikachu-session",
		Value:    id,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   3600, // 1小时过期
	}
	http.SetCookie(w, cookie)
}

func generateSessionID() string {
	return "session_" + time.Now().Format("20060102150405") + "_" + randString(16)
}

func randString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
