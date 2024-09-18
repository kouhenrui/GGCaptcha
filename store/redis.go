package store

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type RediStore struct {
	Store *redis.Client
}

//type RedisOptions struct {
//	Host     string
//	Port     string
//	Db       int
//	PoolSize int
//	MaxRetry int
//}

//	func NewRediStore(option RedisOptions) *RediStore {
//		redisClients := redis.NewClient(&redis.Options{
//			Addr: option.Host + ":" + option.Port,
//			//Username:   redisCon.UserName,
//			//Password:   redisCon.PassWord,
//			DB:         option.Db,
//			PoolSize:   option.PoolSize,
//			MaxRetries: option.MaxRetry,
//		})
//		_, err := redisClients.Ping(context.Background()).Result()
//		if err != nil {
//			log.Printf("redis connect times over %v,please check %v\n", option.MaxRetry, err.Error())
//			panic(fmt.Sprintf("redis connect error,get fauiler %v", err.Error()))
//		}
//		return &RediStore{Store: redisClients}
//	}
func (r *RediStore) Set(id string, value string, t time.Duration) error {
	return r.Store.Set(context.TODO(), id, []byte(value), t).Err()
}

func (r *RediStore) Get(id string, clear bool) (string, error) {
	err := r.Store.Exists(context.TODO(), id).Err()
	if err != nil {
		return "", err
	}
	value, err := r.Store.Get(context.TODO(), id).Result()
	if clear {
		err = r.Store.Del(context.TODO(), id).Err()
		if err != nil {
			return "", err
		}
	}
	return value, nil
}
func (r *RediStore) Verify(id, answer string, clear bool) bool {
	if r.Exist(id) {
		value, err := r.Get(id, clear)
		if err != nil {
			return false
		}
		if answer == value {
			return true
		}
		return false
	}
	return false
}
func (r *RediStore) Exist(id string) bool {
	_, err := r.Store.Exists(context.TODO(), id).Result()
	if err != nil {
		return false
	}
	return true
}
