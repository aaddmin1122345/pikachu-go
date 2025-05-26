package main

import (
	"fmt"
	"net/http"
	"pikachu-go/route"
)

func main() {
	// 设置路由
	route.InitRoutes()
	// 启动服务
	fmt.Println("服务器启动于 http://localhost:8888")
	if err := http.ListenAndServe(":8888", nil); err != nil {
		fmt.Println("服务器启动失败: ", err)
	}
}
