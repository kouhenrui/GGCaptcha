package GGCaptcha

import (
	"GGCaptcha/inter"
	"GGCaptcha/utils"
	"strconv"
	"time"
)

type GGCaptcha struct {
	driver  inter.Driver
	store   inter.Store
	exptime time.Duration
	//windowTime time.Duration
	//limit      int
}

type CaptchaType = string

const (
	StringCaptcha     CaptchaType = "DriverString"
	MathCaptcha       CaptchaType = "DriverMath"
	MathStringCaptcha CaptchaType = "DriverMathString"
	PuzzleCaptcha     CaptchaType = "DriverPuzzle"
)

func NewGGCaptcha(driver inter.Driver, store inter.Store, t time.Duration) *GGCaptcha {
	return &GGCaptcha{
		driver:  driver,
		store:   store,
		exptime: t,
	}
}

func (g *GGCaptcha) GenerateGGCaptcha() (id, content string, err error) {
	content, answer, err := g.driver.GenerateDriverString()
	if err != nil {
		return "", "", err
	}
	id = utils.RandStr(8)
	err = g.store.Set(id, answer, g.exptime)
	if err != nil {
		return "", "", err
	}
	return id, content, nil
}
func (g *GGCaptcha) VerifyGGCaptcha(id, answer string, clear bool) bool {
	return g.store.Verify(id, answer, clear)
}

/*
 * @Title
 * @Description 使用算术验证
 * @Param
 * @return
 * @Author Acer
 * @Date 2024/9/14
 */

func (g *GGCaptcha) GenerateDriverMath() (id, content string, err error) {
	id = utils.RandStr(8)
	content, answer := g.driver.GenerateDriverMath()
	err = g.store.Set(id, answer, g.exptime)
	if err != nil {
		return "", "", err
	}
	return id, content, nil
}

/*
 * @Title
 * @Description 生成计算式图片验证码
 * @Param
 * @return
 * @Author Acer
 * @Date 2024/9/14
 */

func (g *GGCaptcha) GenerateDriverMathString() (id, content string, err error) {
	id = utils.RandStr(8)
	content, answer, err := g.driver.GenerateDriverMathString()
	if err != nil {
		return "", "", err
	}
	err = g.store.Set(id, answer, g.exptime)
	if err != nil {
		return "", "", err
	}
	return id, content, nil
}

/*
 * @Title
 * @Description 滑动验证码
 * @Param
 * @return
 * @Author Acer
 * @Date 2024/9/14
 */

func (g *GGCaptcha) GenerateDriverPuzzle() (id, bgImage, puzzleImage string, err error) {
	id = utils.RandStr(8)
	bgImage, puzzleImage, answer, err := g.driver.GenerateDriverPuzzle()
	if err != nil {
		return "", "", "", err
	}
	err = g.store.Set(id, strconv.Itoa(answer), g.exptime)
	if err != nil {
		return "", "", "", err
	}
	return id, bgImage, puzzleImage, nil
}

/*
 * @Title
 * @Description 刷新接口
 * @Param
 * @return
 * @Author Acer
 * @Date 2024/9/14
 */
func (g *GGCaptcha) RefreshCaptcha(ctype CaptchaType) (id, content string, err error) {
	switch ctype {
	case StringCaptcha:
		id, content, err = g.GenerateGGCaptcha()
	case MathCaptcha:
		id, content, err = g.GenerateDriverMath()
	case MathStringCaptcha:
		id, content, err = g.GenerateDriverMathString()
	}
	return id, content, err
}
