package utils

import (
	"encoding/base64"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"image/color"
	"log"
	"math/rand"
	"os"
	"time"
)

/*
 * @Title
 * @Description 保存base64至本地
 * @Param
 * @return
 * @Author Acer
 * @Date 2024/9/13
 */
func SaveBase64(base, outPath string) error {
	//if strings.Contains(base, "base64,") {
	//	base = strings.Split(base, "base64,")[1]
	//}
	imgData, err := base64.StdEncoding.DecodeString(base)
	if err != nil {
		log.Fatalf("转存失败%s", err)
		return err
	}
	err = os.WriteFile(outPath, imgData, 0644)
	if err != nil {
		log.Fatalf("写入文件失败%s", err)
		return err
	}
	return nil
}

// 加载字体
func LoadDefaultFontFace() font.Face {
	// 使用 Go 自带的字体，也可以使用本地字体文件
	ttf, err := opentype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}
	face, err := opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    20,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	return face
}
func LoadFont(fontPath string, sizePoint float64) font.Face {
	// 读取自定义字体文件
	fontBytes, err := os.ReadFile(fontPath)
	if err != nil {
		log.Fatalf("读取字体文件失败: %v", err)
		//panic(err)
	}

	// 解析字体文件
	ttf, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Fatalf("解析字体文件失败: %v", err)
	}

	// 创建字体 Face
	face := truetype.NewFace(ttf, &truetype.Options{
		Size:    sizePoint,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	return face
}
func RandStr(n int) (randStr string) {
	chars := "ABCDEFGHIJKMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz23456789"
	charsLen := len(chars)
	if n > 10 {
		n = 10
	}
	rand.NewSource(time.Now().UnixNano())
	for i := 0; i < n; i++ {
		randIndex := rand.Intn(charsLen)
		randStr += chars[randIndex : randIndex+1]
	}
	return randStr
}

// 随机颜色
func RandColor(maxColor int) (r, g, b, a int) {
	r = int(uint8(rand.Intn(maxColor)))
	g = int(uint8(rand.Intn(maxColor)))
	b = int(uint8(rand.Intn(maxColor)))
	a = int(uint8(rand.Intn(255)))
	return r, g, b, a
}
func RandColorRGBA(maxColor int) (c color.NRGBA) {
	c.R = uint8(rand.Intn(maxColor))
	c.G = uint8(rand.Intn(maxColor))
	c.B = uint8(rand.Intn(maxColor))
	c.A = uint8(rand.Intn(255))
	return
}

// 随机坐标
func RandPos(width, height int) (x float64, y float64) {
	x = rand.Float64() * float64(width)
	y = rand.Float64() * float64(height)
	return x, y
}

// 随机颜色范围
func RandColorRange(miniColor, maxColor int) (r, g, b, a int) {
	if miniColor > maxColor {
		miniColor = 0
		maxColor = 255
	}
	r = int(uint8(rand.Intn(maxColor-miniColor) + miniColor))
	g = int(uint8(rand.Intn(maxColor-miniColor) + miniColor))
	b = int(uint8(rand.Intn(maxColor-miniColor) + miniColor))
	a = int(uint8(rand.Intn(maxColor-miniColor) + miniColor))
	return r, g, b, a
}
