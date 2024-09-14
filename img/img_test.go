package img

import (
	"GGCaptcha/utils"
	"testing"
)

func TestImg_GenerateDriverString(t *testing.T) {
	img := LoadLocalImg("../utils/dark.png")
	driver := NewDriverString(img)
	content, answer, err := driver.GenerateDriverString()
	if err != nil {
		t.Fatalf("生成图片失败%s", err)
	}
	t.Log(content)

	err = utils.SaveBase64(content, "../1.jpg")
	if err != nil {
		t.Fatalf("保存图片错误%s", err)
	}
	t.Logf(answer)
}

func TestImg_GenerateDriverMath(t *testing.T) {
	driver := NewDriverString()
	content, answer := driver.GenerateDriverMath()

	t.Log(content, answer)
	t.Logf("success")

}

func Test_GenerateDriverMathString(t *testing.T) {
	driver := NewDriverString()
	content, answer, err := driver.GenerateDriverMathString()
	if err != nil {
		t.Fatalf("生成滑动错误%s", err)
	}
	utils.SaveBase64(content, "1.jpg")
	t.Log(answer)
}

func Test_GenerateDriverPuzzle(t *testing.T) {
	driver := NewDriverString()
	bgImage, puzzleImage, targetX, err := driver.GenerateDriverPuzzle()
	if err != nil {
		t.Fatalf("生成滑动错误%s", err)
	}
	utils.SaveBase64(bgImage, "1.jpg")
	utils.SaveBase64(puzzleImage, "2.jpg")
	t.Log(targetX)
}
