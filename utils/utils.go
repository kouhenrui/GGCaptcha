package utils

import (
	"encoding/base64"
	"fmt"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"image/color"
	"log"
	"math/rand"
	"os"
	"strconv"
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

// 随机算术
func RandMath(n int) (expression string, answer string) {
	rand.NewSource(time.Now().UnixNano())
	// 初始化第一个数字
	num := rand.Intn(10) + 1
	expression = fmt.Sprintf("%d", num)
	result := num

	// 定义可用的运算符
	operators := []string{"+", "-", "*", "/"}

	for i := 1; i < n; i++ {
		// 随机选择一个运算符
		operator := operators[rand.Intn(4)]

		// 随机选择一个数字
		var nextNum int
		if operator == "/" {
			// 对于除法，确保能够整除
			nextNum = rand.Intn(9) + 1 // 避免除数为0
			// 确保除法不会导致负结果
			if result < 0 {
				nextNum = rand.Intn(9) + 1
				result *= nextNum
			}
		} else {
			nextNum = rand.Intn(9) + 1 // 避免为0
		}

		// 根据操作符更新结果
		switch operator {
		case "+":
			result += nextNum
		case "-":
			// 确保减法后结果为非负数
			if result-nextNum < 0 {
				// 如果减法会导致负数，重新选择一个操作数
				nextNum = result
			}
			result -= nextNum
		case "*":
			result *= nextNum
		case "/":
			// 仅在除法时更新结果
			// 仅在除法时更新结果
			if result%nextNum != 0 {
				// 确保结果整除
				result += nextNum - (result % nextNum)
			}
			result /= nextNum
		}

		// 构建表达式
		expression += fmt.Sprintf("%s%d", operator, nextNum)
	}

	//// 确保除法后的结果是整数
	//if len(expression) > 0 && expression[len(expression)-1] == '/' {
	//	// 如果最后一个操作符是除法，修正结果
	//	expression = expression[:len(expression)-1] // 移除最后的除法符号
	//}

	expression += "="
	answer = strconv.Itoa(result)
	return expression, answer
}

func RandMath2() (expression string, answer string) {
	rand.NewSource(time.Now().UnixNano())
	// 初始化第一个数字
	num := rand.Intn(20) + 1
	expression = fmt.Sprintf("%d", num)
	result := num
	// 定义可用的运算符
	operators := []string{"+", "-", "*", "/"}
	operator := operators[rand.Intn(len(operators))]
	var nextNum = rand.Intn(20) + 1
	if result < nextNum {
		if operator == "/" {
			// 仅在除法时更新结果
			for result%nextNum != 0 {
				nextNum--
			}
			//if result%nextNum != 0 {
			//	select {
			//		case result  +=1:
			//		if result%nextNum == 0 {
			//			break
			//		}
			//	}
			//	// 确保结果整除
			//	result += nextNum - (result % nextNum)
			//}
			result /= nextNum
		}
	}
	result = nextNum
	// 构建表达式
	expression += fmt.Sprintf("%s%d", operator, nextNum)

	expression += "="
	answer = strconv.Itoa(result)
	return expression, answer
}

// 随机字符
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
