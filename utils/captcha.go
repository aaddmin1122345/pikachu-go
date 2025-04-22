package utils

import (
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"net/http"
	"time"
)

// 生成验证码图像并返回验证码字符串
func GenerateCaptcha(w http.ResponseWriter) string {
	width, height := 80, 30
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 设置背景色为白色
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.RGBA{255, 255, 255, 255})
		}
	}

	// 生成随机字符串
	rand.Seed(time.Now().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	code := ""
	for i := 0; i < 4; i++ {
		idx := rand.Intn(len(chars))
		code += string(chars[idx])
	}

	// 添加字符
	for i, ch := range code {
		// 随机颜色
		r := uint8(rand.Intn(100))
		g := uint8(rand.Intn(100))
		b := uint8(rand.Intn(100))
		col := color.RGBA{r, g, b, 255}

		// 在随机位置绘制字符
		x := 10 + i*15 + rand.Intn(5)
		y := 10 + rand.Intn(10)
		drawChar(img, x, y, string(ch), col)
	}

	// 添加干扰线
	for i := 0; i < 3; i++ {
		r := uint8(rand.Intn(255))
		g := uint8(rand.Intn(255))
		b := uint8(rand.Intn(255))
		col := color.RGBA{r, g, b, 255}

		x1 := rand.Intn(width)
		y1 := rand.Intn(height)
		x2 := rand.Intn(width)
		y2 := rand.Intn(height)
		drawLine(img, x1, y1, x2, y2, col)
	}

	// 添加干扰点
	for i := 0; i < 200; i++ {
		r := uint8(rand.Intn(255))
		g := uint8(rand.Intn(255))
		b := uint8(rand.Intn(255))
		col := color.RGBA{r, g, b, 255}

		x := rand.Intn(width)
		y := rand.Intn(height)
		img.Set(x, y, col)
	}

	// 设置响应头
	w.Header().Set("Content-Type", "image/png")
	// 输出图像
	png.Encode(w, img)

	return code
}

// 绘制字符（简单实现）
func drawChar(img *image.RGBA, x, y int, char string, col color.RGBA) {
	// 简单绘制字符，实际可用更好的字体库
	for i := 0; i < len(char); i++ {
		for j := 0; j < 8; j++ {
			img.Set(x+i, y+j, col)
		}
	}
}

// 绘制线
func drawLine(img *image.RGBA, x1, y1, x2, y2 int, col color.RGBA) {
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	sx, sy := 1, 1
	if x1 >= x2 {
		sx = -1
	}
	if y1 >= y2 {
		sy = -1
	}
	err := dx - dy

	for {
		img.Set(x1, y1, col)
		if x1 == x2 && y1 == y2 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
