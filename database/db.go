package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// 数据库配置
var (
	DBHost     = getEnv("DB_HOST", "localhost")
	DBPort     = getEnv("DB_PORT", "5432")
	DBUser     = getEnv("DB_USER", "pgsql")
	DBPassword = getEnv("DB_PASSWORD", "pgsql")
	DBName     = getEnv("DB_NAME", "pikachu-go")
	DBSSLMode  = getEnv("DB_SSLMODE", "disable")
)

// DB 是数据库的全局变量，表示数据库连接池
var DB *sql.DB

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

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
	
	-- 插入一些测试数据，如果用户表为空
	INSERT INTO users (username, password, email)
	SELECT 'pikachu', 'pikachu123', 'pikachu@example.com'
	WHERE NOT EXISTS (SELECT 1 FROM users WHERE username = 'pikachu');
	
	INSERT INTO users (username, password, email)
	SELECT 'admin', 'admin123', 'admin@example.com'
	WHERE NOT EXISTS (SELECT 1 FROM users WHERE username = 'admin');
	`)

	if err != nil {
		log.Printf("数据库表格创建失败：%v", err)
		return
	}

	log.Println("数据库初始化成功")
}
