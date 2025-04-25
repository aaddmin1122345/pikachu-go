package utils

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
	"net/http"
	"time"
)

// 生成验证码图像并返回验证码字符串
func GenerateCaptcha(w http.ResponseWriter) string {
	// 确保随机数是真随机的
	rand.Seed(time.Now().UnixNano())

	// 增加验证码图片尺寸
	width, height := 120, 40
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 设置背景色为白色
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// 生成随机字符串
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	code := ""
	for i := 0; i < 4; i++ {
		idx := rand.Intn(len(chars))
		code += string(chars[idx])
	}

	// 添加字符
	for i, ch := range code {
		// 随机颜色
		textColor := color.RGBA{
			R: uint8(rand.Intn(100)),
			G: uint8(rand.Intn(100)),
			B: uint8(rand.Intn(100)),
			A: 255,
		}

		// 调整字符位置和大小
		x := 20 + i*25 // 增加字符间距
		y := 10        // 调整垂直位置
		drawBoldChar(img, x, y, string(ch), textColor)
	}

	// 添加干扰线
	for i := 0; i < 4; i++ {
		lineColor := color.RGBA{
			R: uint8(rand.Intn(255)),
			G: uint8(rand.Intn(255)),
			B: uint8(rand.Intn(255)),
			A: 255,
		}

		x1 := rand.Intn(width)
		y1 := rand.Intn(height)
		x2 := rand.Intn(width)
		y2 := rand.Intn(height)
		drawLine(img, x1, y1, x2, y2, lineColor)
	}

	// 添加干扰点
	for i := 0; i < 200; i++ {
		noiseColor := color.RGBA{
			R: uint8(rand.Intn(255)),
			G: uint8(rand.Intn(255)),
			B: uint8(rand.Intn(255)),
			A: 255,
		}

		x := rand.Intn(width)
		y := rand.Intn(height)
		img.Set(x, y, noiseColor)
	}

	// 设置响应头
	w.Header().Set("Content-Type", "image/png")
	// 输出图像
	png.Encode(w, img)

	return code
}

// 绘制粗体字符（简单但更清晰的实现，增加字体大小）
func drawBoldChar(img *image.RGBA, x, y int, char string, col color.RGBA) {
	// 根据字符选择一个简单的点阵字体模式
	patterns := map[rune][]string{
		'a': {
			"   XXXX   ",
			"  XX  XX  ",
			" XX    XX ",
			" XXXXXXXX ",
			" XX    XX ",
			" XX    XX ",
			" XX    XX ",
		},
		'b': {
			" XXXXXXX  ",
			" XX    XX ",
			" XX    XX ",
			" XXXXXXX  ",
			" XX    XX ",
			" XX    XX ",
			" XXXXXXX  ",
		},
		'c': {
			"  XXXXXX  ",
			" XX    XX ",
			" XX       ",
			" XX       ",
			" XX       ",
			" XX    XX ",
			"  XXXXXX  ",
		},
		'd': {
			" XXXXXXX  ",
			" XX    XX ",
			" XX    XX ",
			" XX    XX ",
			" XX    XX ",
			" XX    XX ",
			" XXXXXXX  ",
		},
		'e': {
			" XXXXXXXX ",
			" XX       ",
			" XX       ",
			" XXXXXX   ",
			" XX       ",
			" XX       ",
			" XXXXXXXX ",
		},
		'f': {
			" XXXXXXXX ",
			" XX       ",
			" XX       ",
			" XXXXXX   ",
			" XX       ",
			" XX       ",
			" XX       ",
		},
		'g': {
			"  XXXXXX  ",
			" XX    XX ",
			" XX       ",
			" XX  XXXX ",
			" XX    XX ",
			" XX    XX ",
			"  XXXXXX  ",
		},
		'h': {
			" XX    XX ",
			" XX    XX ",
			" XX    XX ",
			" XXXXXXXX ",
			" XX    XX ",
			" XX    XX ",
			" XX    XX ",
		},
		'i': {
			" XXXXXXXX ",
			"    XX    ",
			"    XX    ",
			"    XX    ",
			"    XX    ",
			"    XX    ",
			" XXXXXXXX ",
		},
		'j': {
			" XXXXXXXX ",
			"     XX   ",
			"     XX   ",
			"     XX   ",
			"     XX   ",
			" XX  XX   ",
			"  XXXX    ",
		},
		'k': {
			" XX    XX ",
			" XX   XX  ",
			" XX  XX   ",
			" XXXXX    ",
			" XX  XX   ",
			" XX   XX  ",
			" XX    XX ",
		},
		'l': {
			" XX       ",
			" XX       ",
			" XX       ",
			" XX       ",
			" XX       ",
			" XX       ",
			" XXXXXXXX ",
		},
		'm': {
			" XX    XX ",
			" XXX  XXX ",
			" XXXXXXXX ",
			" XX XX XX ",
			" XX    XX ",
			" XX    XX ",
			" XX    XX ",
		},
		'n': {
			" XX    XX ",
			" XXX   XX ",
			" XXXX  XX ",
			" XX XX XX ",
			" XX  XXXX ",
			" XX   XXX ",
			" XX    XX ",
		},
		'o': {
			"  XXXXXX  ",
			" XX    XX ",
			" XX    XX ",
			" XX    XX ",
			" XX    XX ",
			" XX    XX ",
			"  XXXXXX  ",
		},
		'p': {
			" XXXXXXX  ",
			" XX    XX ",
			" XX    XX ",
			" XXXXXXX  ",
			" XX       ",
			" XX       ",
			" XX       ",
		},
		'q': {
			"  XXXXXX  ",
			" XX    XX ",
			" XX    XX ",
			" XX    XX ",
			" XX  X XX ",
			" XX   XX  ",
			"  XXXX XX ",
		},
		'r': {
			" XXXXXXX  ",
			" XX    XX ",
			" XX    XX ",
			" XXXXXXX  ",
			" XX  XX   ",
			" XX   XX  ",
			" XX    XX ",
		},
		's': {
			"  XXXXXX  ",
			" XX    XX ",
			" XX       ",
			"  XXXXXX  ",
			"       XX ",
			" XX    XX ",
			"  XXXXXX  ",
		},
		't': {
			" XXXXXXXX ",
			"    XX    ",
			"    XX    ",
			"    XX    ",
			"    XX    ",
			"    XX    ",
			"    XX    ",
		},
		'u': {
			" XX    XX ",
			" XX    XX ",
			" XX    XX ",
			" XX    XX ",
			" XX    XX ",
			" XX    XX ",
			"  XXXXXX  ",
		},
		'v': {
			" XX    XX ",
			" XX    XX ",
			" XX    XX ",
			" XX    XX ",
			" XX    XX ",
			"  XX  XX  ",
			"   XXXX   ",
		},
		'w': {
			" XX    XX ",
			" XX    XX ",
			" XX    XX ",
			" XX XX XX ",
			" XXXXXXXX ",
			" XXX  XXX ",
			" XX    XX ",
		},
		'x': {
			" XX    XX ",
			" XX    XX ",
			"  XX  XX  ",
			"   XXXX   ",
			"  XX  XX  ",
			" XX    XX ",
			" XX    XX ",
		},
		'y': {
			" XX    XX ",
			" XX    XX ",
			"  XX  XX  ",
			"   XXXX   ",
			"    XX    ",
			"    XX    ",
			"    XX    ",
		},
		'z': {
			" XXXXXXXX ",
			"      XX  ",
			"     XX   ",
			"    XX    ",
			"   XX     ",
			"  XX      ",
			" XXXXXXXX ",
		},
		'0': {
			"  XXXXXX  ",
			" XX    XX ",
			" XX   XXX ",
			" XX  XXXX ",
			" XXXX  XX ",
			" XXX   XX ",
			"  XXXXXX  ",
		},
		'1': {
			"    XX    ",
			"   XXX    ",
			"    XX    ",
			"    XX    ",
			"    XX    ",
			"    XX    ",
			" XXXXXXXX ",
		},
		'2': {
			"  XXXXXX  ",
			" XX    XX ",
			"      XX  ",
			"    XXX   ",
			"  XX      ",
			" XX       ",
			" XXXXXXXX ",
		},
		'3': {
			" XXXXXXXX ",
			"       XX ",
			"       XX ",
			"   XXXXX  ",
			"       XX ",
			" XX    XX ",
			"  XXXXXX  ",
		},
		'4': {
			" XX    XX ",
			" XX    XX ",
			" XX    XX ",
			" XXXXXXXX ",
			"       XX ",
			"       XX ",
			"       XX ",
		},
		'5': {
			" XXXXXXXX ",
			" XX       ",
			" XX       ",
			" XXXXXXX  ",
			"       XX ",
			" XX    XX ",
			"  XXXXXX  ",
		},
		'6': {
			"  XXXXXX  ",
			" XX       ",
			" XX       ",
			" XXXXXXX  ",
			" XX    XX ",
			" XX    XX ",
			"  XXXXXX  ",
		},
		'7': {
			" XXXXXXXX ",
			"       XX ",
			"      XX  ",
			"     XX   ",
			"    XX    ",
			"   XX     ",
			"  XX      ",
		},
		'8': {
			"  XXXXXX  ",
			" XX    XX ",
			" XX    XX ",
			"  XXXXXX  ",
			" XX    XX ",
			" XX    XX ",
			"  XXXXXX  ",
		},
		'9': {
			"  XXXXXX  ",
			" XX    XX ",
			" XX    XX ",
			"  XXXXXXX ",
			"       XX ",
			"       XX ",
			"  XXXXXX  ",
		},
	}

	// 获取字符的点阵
	if pattern, ok := patterns[rune(char[0])]; ok {
		for i, line := range pattern {
			for j, pixel := range line {
				if pixel == 'X' {
					// 画点块
					img.Set(x+j, y+i, col)
				}
			}
		}
	} else {
		// 如果字符没有预定义的点阵，就画一个简单的方块
		for i := 0; i < 7; i++ {
			for j := 0; j < 10; j++ {
				img.Set(x+j, y+i, col)
			}
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
