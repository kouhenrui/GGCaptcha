package GGCaptcha

import (
	"GGCaptcha/img"
	"GGCaptcha/store"
	"testing"
	"time"
)

func TestNewGGCaptcha(t *testing.T) {
	driver := img.NewDriverString()
	redisOption := store.RedisOptions{
		Host:     "192.168.245.22",
		Port:     "6379",
		Db:       4,
		PoolSize: 10,
		MaxRetry: 5,
	}
	redisStore := store.NewRediStore(redisOption)
	//localStore := store.NewLocalStore()
	ggcaptcha := NewGGCaptcha(driver, redisStore, 1*time.Minute)
	id, content, err := ggcaptcha.GenerateGGCaptcha()
	if err != nil {
		t.Fatalf("生成文件错误%s", err)
	}
	t.Log(id, content)
}

func Test_GenerateDriverMath(t *testing.T) {
	localStore := store.NewLocalStore()
	driver := img.NewDriverString()
	ggcaptcha := NewGGCaptcha(driver, localStore, time.Minute)

	id, content, err := ggcaptcha.GenerateDriverMath()
	if err != nil {
		t.Fatalf("生成算术验证码错误%s", err)
	}
	t.Log(id, content)
}

func TestGGCaptcha_GenerateDriverMathString(t *testing.T) {
	localStore := store.NewLocalStore()
	driver := img.NewDriverString()
	ggcaptcha := NewGGCaptcha(driver, localStore, time.Minute)
	id, content, err := ggcaptcha.GenerateDriverMathString()
	if err != nil {
		t.Fatalf("生成算术验证码错误%s", err)
	}
	t.Log(id, content)
}

func Test_GenerateDriverPuzzle(t *testing.T) {
	localStore := store.NewLocalStore()
	driver := img.NewDriverString()
	ggcaptcha := NewGGCaptcha(driver, localStore, time.Minute)
	id, bgImage, puzzleImage, err := ggcaptcha.GenerateDriverPuzzle()
	if err != nil {
		t.Fatalf("生成滑动验证码错误%s", err)
	}
	t.Log(id, bgImage, puzzleImage)
}

func Test_RefreshCaptcha(t *testing.T) {
	localStore := store.NewLocalStore()
	driver := img.NewDriverString()
	ggcaptcha := NewGGCaptcha(driver, localStore, time.Minute)
	id, content, err := ggcaptcha.RefreshCaptcha(StringCaptcha)
	if err != nil {
		t.Fatalf("刷新验证码错误%s", err)
	}
	t.Log(id, content)
}
