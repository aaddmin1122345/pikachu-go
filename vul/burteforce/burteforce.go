package burteforce

import (
	"net/http"

	"pikachu-go/templates"
)

// BurteforceHandler 渲染暴力破解模块概述页
func BurteforceHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 设置导航栏的激活状态
		data := templates.NewPageData2(1, 2, "")
		renderer.RenderPage(w, "burteforce/burteforce.html", data)
	}
}
