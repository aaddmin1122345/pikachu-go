package burteforce

import (
	"net/http"
	"pikachu-go/templates"
)

// BurteforceHandler 渲染爆破攻击概述页
func BurteforceHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewPageData2(12, 13, "")
		renderer.RenderPage(w, "burteforce/burteforce.html", data)
	}
}
