package ggs

import (
	"GGCaptcha/img"
	"GGCaptcha/inter"
	"GGCaptcha/store"
	"GGCaptcha/utils"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
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
		i.Height = 60
		i.Width = 120
		i.Count = 4
		i.SizePoint = 20
		i.NoiseCount = 4
		i.Source = "ABCDEFGHIJKMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz0123456789"
		i.SourceLength = len(i.Source)
		i.Bgcolor = utils.RandColorRGBA(255)
		i.FontStyle = utils.LoadFontFace(i.SizePoint)
	} else {
		i = imgOptions[0]
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
