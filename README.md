# GGCaptcha
使用gg图像库和Redis以及本地缓存封装的验证码组件,包含简单数式计算验证、图片字符验证、图片数式计算验证和滑动图片验证。
支持使用自定义图片或使用本组件默认图片参数
```
	localStore := GGCaptcha.NewLocalStore()
	redisOption := GGCaptcha.RedisOptions{
		Host:     "192.168.245.22",
		Port:     "6379",
		Db:       4,
		PoolSize: 10,
		MaxRetry: 5,
	}
use default captcha for example(使用默认参数示例)
	driver := GGCaptcha.NewDriverString()
	redisStore := GGCaptcha.NewRediStore(redisOption)
	ggcaptcha := GGCaptcha.NewGGCaptcha(driver, redisStore)
	id, content, err := ggcaptcha.GenerateGGCaptcha()

use img as background
	img := GGCaptcha.LoadLocalImg("../utils/dark.png")
	driver := GGCaptcha.NewDriverString(img)
	content, answer, err := GGCaptcha.GenerateDriverString()
```
