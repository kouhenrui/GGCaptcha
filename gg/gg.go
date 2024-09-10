package ggs

import (
	"GGCaptcha/inter"
	"GGCaptcha/utils"
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
