package database

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// 数据库配置
var (
	DBHost     = "localhost"
	DBPort     = "5432"
	DBUser     = "pgsql"
	DBPassword = "pgsql"
	DBName     = "pikachu-go"
	DBSSLMode  = "disable"
)

// DB 是数据库的全局变量，表示数据库连接池
var DB *sql.DB

// getEnv 获取环境变量，如果不存在则返回默认值
// func getEnv(key, defaultValue string) string {
// 	if value, exists := os.LookupEnv(key); exists {
// 		return value
// 	}
// 	return defaultValue
// }

// GetDSN 获取PostgreSQL数据库连接字符串
func GetDSN() string {
	return "host=" + DBHost + " port=" + DBPort + " user=" + DBUser +
		" password=" + DBPassword + " dbname=" + DBName + " sslmode=" + DBSSLMode
}

// InitializeDatabase 初始化数据库，连接PostgreSQL并创建表格
func InitializeDatabase() {
	var err error

	// 连接PostgreSQL数据库
	dsn := GetDSN()
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Printf("数据库连接失败：%v", err)
		return
	}

	// 测试连接
	if err := DB.Ping(); err != nil {
		log.Printf("数据库ping失败：%v", err)
		return
	}

	// 初始化数据库表格 - 使用PostgreSQL语法
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		email TEXT
	);
	CREATE TABLE IF NOT EXISTS sessions (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id),
		token TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS logins (
		id SERIAL PRIMARY KEY,
		username TEXT,
		password TEXT,
		ip_address TEXT,
		login_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS message (
		id SERIAL PRIMARY KEY,
		content TEXT,
		time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS xssblind (
		id SERIAL PRIMARY KEY,
		time TIMESTAMP,
		content TEXT,
		name TEXT
	);
	CREATE TABLE IF NOT EXISTS member (
		id SERIAL PRIMARY KEY,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		sex TEXT,
		phonenum TEXT,
		address TEXT,
		email TEXT
	);
	`)

	if err != nil {
		log.Printf("数据库表格创建失败：%v", err)
		return
	}

	// MD5加密密码
	pikachuPasswordMD5 := fmt.Sprintf("%x", md5.Sum([]byte("pikachu123")))
	adminPasswordMD5 := fmt.Sprintf("%x", md5.Sum([]byte("admin123")))
	memberPasswordMD5 := fmt.Sprintf("%x", md5.Sum([]byte("123456")))

	// 插入pikachu用户
	_, err = DB.Exec(`
	INSERT INTO users (username, password, email)
	SELECT 'pikachu', $1, 'pikachu@example.com'
	WHERE NOT EXISTS (SELECT 1 FROM users WHERE username = 'pikachu')
	`, pikachuPasswordMD5)

	if err != nil {
		log.Printf("pikachu用户创建失败：%v", err)
	}

	// 插入admin用户
	_, err = DB.Exec(`
	INSERT INTO users (username, password, email)
	SELECT 'admin', $1, 'admin@example.com'
	WHERE NOT EXISTS (SELECT 1 FROM users WHERE username = 'admin')
	`, adminPasswordMD5)

	if err != nil {
		log.Printf("admin用户创建失败：%v", err)
	}

	// 插入member表数据，使用参数化查询
	memberData := []struct {
		username string
		sex      string
		phonenum string
		address  string
		email    string
	}{
		{"vince", "男", "18612345678", "北京市", "vince@example.com"},
		{"allen", "男", "18712345678", "上海市", "allen@example.com"},
		{"kobe", "男", "18812345678", "广州市", "kobe@example.com"},
		{"grady", "男", "18912345678", "深圳市", "grady@example.com"},
		{"kevin", "男", "13012345678", "杭州市", "kevin@example.com"},
		{"lucy", "女", "13112345678", "南京市", "lucy@example.com"},
		{"lili", "女", "13212345678", "成都市", "lili@example.com"},
	}

	for _, member := range memberData {
		_, err = DB.Exec(`
		INSERT INTO member (username, password, sex, phonenum, address, email)
		SELECT $1, $2, $3, $4, $5, $6
		WHERE NOT EXISTS (SELECT 1 FROM member WHERE username = $1)
		`, member.username, memberPasswordMD5, member.sex, member.phonenum, member.address, member.email)

		if err != nil {
			log.Printf("创建会员 %s 失败：%v", member.username, err)
		}
	}

	log.Println("数据库初始化成功")
}
