package store

import (
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"
)

type LocalStore struct {
	item sync.Map
	lock sync.Mutex
}
type cache struct {
	Value string
	T     int64
}

func newcache(value string, t time.Duration) *cache {
	rand.NewSource(time.Now().UnixNano())
	return &cache{
		Value: value,
		T:     time.Now().Add(t).UnixNano(),
	}
}
func NewLocalStore() *LocalStore {
	return &LocalStore{
		item: sync.Map{},
		lock: sync.Mutex{},
	}
}
func (l *LocalStore) Set(id string, value string, t time.Duration) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	cach := newcache(value, t)
	l.item.Store(id, cach)
	return nil
}
func (l *LocalStore) Exist(id string) bool {
	_, ok := l.item.Load(id)
	return ok
}
func (l *LocalStore) Get(id string, clear bool) (string, error) {
	if l.Exist(id) {
		value, _ := l.item.Load(id)
		cacheItem := value.(*cache)
		// 检查是否过期
		if time.Now().UnixNano() > cacheItem.T || clear {
			l.item.Delete(id)
			return "", errors.New("缓存已过期")
		}
		return cacheItem.Value, nil
	} else {
		return "", nil
	}
}
func (l *LocalStore) Verify(id, answer string, clear bool) bool {
	value, err := l.Get(id, clear)
	log.Println(value)
	if err != nil {
		return false
	}
	if answer == value {
		return true
	}
	return false

}
