package database

import (
	"crypto/md5"
	"fmt"
	"log"
)

// UpdateExistingPasswords 更新数据库中的明文密码为MD5格式
func UpdateExistingPasswords() {
	if DB == nil {
		log.Println("数据库未初始化，无法更新密码")
		return
	}

	// 查询所有用户
	rows, err := DB.Query("SELECT id, username, password FROM users")
	if err != nil {
		log.Printf("查询用户失败：%v", err)
		return
	}
	defer rows.Close()

	// 对每个用户，检查密码是否已经是MD5格式（32位十六进制字符串）
	for rows.Next() {
		var id int
		var username, password string
		if err := rows.Scan(&id, &username, &password); err != nil {
			log.Printf("读取用户数据失败：%v", err)
			continue
		}

		// 检查密码是否可能是明文（简单判断：长度小于32且不全是十六进制字符）
		if len(password) != 32 || !isMD5Format(password) {
			// 将明文密码转换为MD5
			md5Password := fmt.Sprintf("%x", md5.Sum([]byte(password)))

			// 更新数据库
			_, err := DB.Exec("UPDATE users SET password = $1 WHERE id = $2", md5Password, id)
			if err != nil {
				log.Printf("更新用户 %s 的密码失败：%v", username, err)
			} else {
				log.Printf("已将用户 %s 的密码更新为MD5格式", username)
			}
		}
	}

	if err := rows.Err(); err != nil {
		log.Printf("遍历用户数据时出错：%v", err)
	}

	log.Println("密码格式更新完成")
}

// isMD5Format 检查字符串是否符合MD5格式（32位十六进制字符）
func isMD5Format(s string) bool {
	if len(s) != 32 {
		return false
	}
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}
