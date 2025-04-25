package burteforce

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
	"net/http"
	"time"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/gorilla/sessions"
	"golang.org/x/image/font/gofont/goregular"
)

var sessionStore = sessions.NewCookieStore([]byte("pikachu-secret-key"))

func init() {
	rand.Seed(time.Now().UnixNano())
}

// VcodeHandler 生成图形验证码并写入 session
func VcodeHandler(w http.ResponseWriter, r *http.Request) {
	// 随机生成验证码
	const chars = "0123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz"
	code := make([]byte, 5)
	for i := range code {
		code[i] = chars[rand.Intn(len(chars))]
	}
	vcode := string(code)

	// 存入 session
	sess, _ := sessionStore.Get(r, "pikachu-session")
	sess.Values["vcode"] = vcode
	_ = sess.Save(r, w)

	// 生成图片
	const width, height = 100, 40
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// 加载字体
	font, _ := truetype.Parse(goregular.TTF)
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(font)
	c.SetFontSize(20)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.NewUniform(color.RGBA{0, 0, 255, 255}))

	// 写入文字
	for i, ch := range vcode {
		pt := freetype.Pt(10+i*18, 28)
		c.DrawString(string(ch), pt)
	}

	w.Header().Set("Content-Type", "image/png")
	_ = png.Encode(w, img)
}
