package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// DB 是数据库的全局变量，表示数据库连接池
var DB *sql.DB

// InitializeDatabase 初始化数据库，连接 PostgreSQL 并创建表格
func InitializeDatabase() {
	var err error

	// 连接 PostgreSQL 数据库
	DB, err = sql.Open("postgres", "user=pgsql password=pgsql dbname=pikachu-go sslmode=disable")
	if err != nil {
		log.Fatal("数据库连接失败：", err)
	}

	// 测试连接
	if err := DB.Ping(); err != nil {
		log.Fatal("数据库连接失败：", err)
	}

	// 初始化数据库表格
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT NOT NULL,
		password TEXT NOT NULL
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
	`)
	if err != nil {
		log.Fatal("数据库表格创建失败：", err)
	}
}
