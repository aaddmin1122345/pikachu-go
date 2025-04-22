package unserilization

import (
	"net/http"
	"pikachu-go/templates"
)

// UnserilizationHandler 渲染漏洞概述页面
func UnserilizationHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewPageData2(90, 91, "")
		renderer.RenderPage(w, "unserilization/unserilization.html", data)
	}
}
