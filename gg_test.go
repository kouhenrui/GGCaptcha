package GGCaptcha

import (
	"github.com/kouhenrui/GGCaptcha/inter"
	"testing"
	"time"
)

func TestUploadImg(t *testing.T) {
	driver, err := LoadLocalImg("dark.png")
	if err != nil {
		t.Error(err)
	}
	localStore := NewLocalStore()

	ggcaptcha := NewGGCaptcha(driver, localStore, 1*time.Minute, 10*time.Minute, 50)
	id, content, err := ggcaptcha.GenerateGGCaptcha()

	if err != nil {
		t.Fatalf("生成文件错误%s", err)
	}
	t.Log(id, content)
}

func TestNewGGCaptcha(t *testing.T) {
	driver := NewDriverString()

	//UserName := "root"
	//Password := "123456"
	//redisOption := RedisOptions{
	//	Host: "192.168.255.2", //"r-bp1dh88jq0pzv6tudo.redis.rds.aliyuncs.com",
	//	Port: "6379",
	//	//UserName: &UserName,
	//	Password: &Password,
	//	Db:       4,
	//	PoolSize: 10,
	//	MaxRetry: 5,
	//}
	//redisStore := NewRediStore(redisOption)
	localStore := NewLocalStore()
	ggcaptcha := NewGGCaptcha(driver, localStore, 1*time.Minute, 10*time.Minute, 50)
	id, content, err := ggcaptcha.GenerateGGCaptcha()
	if err != nil {
		t.Fatalf("生成文件错误%s", err)
	}
	t.Log(id, content)
}

func Test_GenerateDriverMath(t *testing.T) {
	localStore := NewLocalStore()
	driver := NewDriverString()
	ggcaptcha := NewGGCaptcha(driver, localStore, 1*time.Minute, 10*time.Minute, 50)

	id, content, err := ggcaptcha.GenerateDriverMath()
	if err != nil {
		t.Fatalf("生成算术验证码错误%s", err)
	}
	//time.Sleep(5 * time.Second)
	//time.Sleep(1 * time.Minute)
	//verify := "25"
	//verifyResult, err := ggcaptcha.VerifyGGCaptcha(id, verify, true)
	//if err != nil {
	//	t.Error("验证码验证失败", err.Error())
	//}
	//t.Log(verifyResult)
	t.Log(id, content)
}

func TestGGCaptcha_GenerateDriverMathString(t *testing.T) {
	localStore := NewLocalStore()
	driver := NewDriverString()
	ggcaptcha := NewGGCaptcha(driver, localStore, 1*time.Minute, 10*time.Minute, 50)
	id, content, err := ggcaptcha.GenerateDriverMathString()
	if err != nil {
		t.Fatalf("生成算术验证码错误%s", err)
	}
	t.Log(id, content)
}

func Test_GenerateDriverPuzzle(t *testing.T) {
	localStore := NewLocalStore()
	driver := NewDriverString()
	ggcaptcha := NewGGCaptcha(driver, localStore, 1*time.Minute, 10*time.Minute, 50)
	id, bgImage, puzzleImage, err := ggcaptcha.GenerateDriverPuzzle()
	if err != nil {
		t.Fatalf("生成滑动验证码错误%s", err)
	}
	t.Log(id, bgImage, puzzleImage)
}

func Test_RefreshCaptcha(t *testing.T) {
	localStore := NewLocalStore()
	driver := NewDriverString()
	ggcaptcha := NewGGCaptcha(driver, localStore, 1*time.Minute, 10*time.Minute, 50)
	id, content, err := ggcaptcha.RefreshCaptcha(inter.StringCaptcha)
	if err != nil {
		t.Fatalf("刷新验证码错误%s", err)
	}
	t.Log(id, content)
}
