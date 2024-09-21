package inter

import (
	"github.com/kouhenrui/GGCaptcha"
	"time"
)

type CaptchaType = GGCaptcha.CaptchaType

type GGCaptchator interface {
	GenerateGGCaptcha() (id, content string, err error)

	GenerateDriverMath() (id, content string, err error)

	GenerateDriverMathString() (id, content string, err error)

	GenerateDriverPuzzle() (id, bgImage, puzzleImage string, err error)

	VerifyGGCaptcha(id, answer string, clear bool) bool

	RefreshCaptcha(ctype CaptchaType) (id, content string, err error)
}

/*
 * @Title
 * @Description 缓存图片存储机制
 * @Author Acer
 * @Date 2024/9/21
 */

type Store interface {
	Set(id string, value string, t time.Duration) error
	Exist(id string) bool
	Get(id string, clear bool) (string, error)
	Verify(id, answer string, clear bool) bool
}

/*
 * @Title
 * @Description 缓存图片验证码驱动
 * @Author Acer
 * @Date 2024/9/21
 */

type Driver interface {
	GenerateDriverString() (content, answer string, err error)
	GenerateDriverMathString() (content, answer string, err error)
	GenerateDriverMath() (content, answer string)
	GenerateDriverPuzzle() (bgImage string, puzzleImage string, targetX int, err error)
}
