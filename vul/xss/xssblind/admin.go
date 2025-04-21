package xssblind

import (
	"net/http"
	"pikachu-go/templates"
	"sync"
)

var (
	stolenData []string
	lock       sync.Mutex
)

func AdminHandler(renderer templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			r.ParseForm()
			content := r.FormValue("data")
			lock.Lock()
			stolenData = append(stolenData, content)
			lock.Unlock()
			w.Write([]byte("ok"))
			return
		}

		lock.Lock()
		html := ""
		for _, d := range stolenData {
			html += d + "<br>"
		}
		lock.Unlock()

		data := templates.NewPageData2(7, 20, html)
		renderer.RenderPage(w, "admin.html", data)
	}
}
