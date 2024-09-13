package ggs

import (
	"GGCaptcha/img"
	"GGCaptcha/inter"
	"GGCaptcha/store"
	"GGCaptcha/utils"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"image/color"
	"log"
	"sync"
	"time"
)

type GGCaptcha struct {
	driver inter.Driver
	store  inter.Store
}

func NewGGCaptcha(driver inter.Driver, store inter.Store) *GGCaptcha {
	return &GGCaptcha{
		driver: driver,
		store:  store,
	}
}

func NewDriverString(imgOptions ...img.Img) *img.Img {
	var i img.Img
	if len(imgOptions) < 1 {
		i = img.DefaultImg()
	} else {
		// 使用传入的自定义配置
		i = imgOptions[0]
		//log.Println(i, "打印参数")
		// 如果没有设置某些字段，则为其提供默认值
		if i.Height == 0 {
			i.Height = 60
		}
		if i.Width == 0 {
			i.Width = 120
		}
		if i.Count == 0 {
			i.Count = 4
		}
		if i.SizePoint == 0 {
			i.SizePoint = 20
		}
		if i.NoiseCount == 0 {
			i.NoiseCount = 4
		}
		if i.Source == "" {
			i.Source = "ABCDEFGHIJKMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz0123456789"
			i.SourceLength = len(i.Source)
		}
		if i.BgColor == (color.NRGBA{}) {
			i.BgColor = utils.RandColorRGBA(255)
		}
		if i.FontStyle == nil {
			i.FontStyle = utils.LoadDefaultFontFace()
		}
		if i.FontColor == nil {
			i.FontColor = utils.RandColorRGBA(255)
		}
	}

	return &i
}
func NewRediStore(option store.RedisOptions) *store.RediStore {
	redisClients := redis.NewClient(&redis.Options{
		Addr: option.Host + ":" + option.Port,
		//Username:   redisCon.UserName,
		//Password:   redisCon.PassWord,
		DB:         option.Db,
		PoolSize:   option.PoolSize,
		MaxRetries: option.MaxRetry,
	})
	_, err := redisClients.Ping(context.Background()).Result()
	if err != nil {
		log.Printf("redis connect times over %v,please check %v\n", option.MaxRetry, err.Error())
		panic(fmt.Sprintf("redis connect error,get fauiler %v", err.Error()))
	}
	return &store.RediStore{Store: redisClients}
}

func NewLocalStore() *store.LocalStore {
	return &store.LocalStore{
		Item: sync.Map{},
		Lock: sync.Mutex{},
	}
}
func (g *GGCaptcha) GenerateGGCaptcha() (id, content, answer string, err error) {
	content, answer, err = g.driver.GenerateDriverString()
	if err != nil {
		return "", "", "", err
	}
	id = utils.RandStr(8)
	err = g.store.Set(id, answer, time.Minute)
	if err != nil {
		return "", "", "", err
	}
	return id, content, answer, nil
}
func (g *GGCaptcha) VerifyGGCaptcha(id, answer string, clear bool) bool {
	return g.store.Verify(id, answer, clear)
}

func UploadImage() {

}
