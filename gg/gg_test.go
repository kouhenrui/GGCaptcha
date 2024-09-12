package ggs

import (
	"testing"
)

func TestNewGGCaptcha(t *testing.T) {
	driver := NewDriverString()
	//redisOption := store.RedisOptions{
	//	Host:     "192.168.245.22",
	//	Port:     "6379",
	//	Db:       4,
	//	PoolSize: 10,
	//	MaxRetry: 5,
	//}
	//redisStore := store.NewRediStore(redisOption)
	localStore := NewLocalStore()
	ggcaptcha := NewGGCaptcha(driver, localStore)
	id, _, answer, err := ggcaptcha.GenerateGGCaptcha()
	if err != nil {
		t.Fatalf("生成文件错误%s", err)
	}
	t.Log(answer)
	t.Log("打印结果id", id)
	b := ggcaptcha.VerifyGGCaptcha(id, answer, false)
	t.Log("打印验证码是否正确", b)
}
