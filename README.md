# GGCaptcha
使用gg图像库和Redis以及本地缓存封装的验证码组件,包含简单数式计算验证、图片字符验证、图片数式计算验证和滑动图片验证。
支持使用自定义图片或使用本组件默认图片参数
```
	driver := GGCaptcha.NewDriverString()
	localStore := GGCaptcha.NewLocalStore()
	redisOption := GGCaptcha.RedisOptions{
		Host:     "192.168.245.22",
		Port:     "6379",
		Db:       4,
		PoolSize: 10,
		MaxRetry: 5,
	}
use default captcha for example(使用默认参数示例)
	redisStore := GGCaptcha.NewRediStore(redisOption)
	ggcaptcha := GGCaptcha.NewGGCaptcha(driver, redisStore)
	id, content, err := ggcaptcha.GenerateGGCaptcha()

use img as background(使用自定义背景图片)
	img := GGCaptcha.LoadLocalImg("../utils/dark.png")
	content, answer, err := GGCaptcha.GenerateDriverString()

use string verify(使用算数验证，直接返回字符类型算术式)
	ggcaptcha := GGCaptcha.NewGGCaptcha(driver, localstore, time.Minute)
	id, content, err := ggcaptcha.GenerateGGCaptcha()

use  string math(使用算术式图片验证，返回图片)
	ggcaptcha := GGCaptcha.NewGGCaptcha(driver, localstore, time.Minute)
	id, content, err := ggcaptcha.GenerateDriverMathString()

use puzzle verify(使用滑动验证码验证)
	ggcaptcha := GGCaptcha.NewGGCaptcha(driver, localstore, time.Minute)
	id, bgImage, puzzleImage, err := ggcaptcha.GenerateDriverPuzzle()

use verify_function verify code(验证密钥是否正确，是否清除)
	ggcaptcha.VerifyGGCaptcha(id,answer,true)


使用代码生成器
	table := Table{
		Name: "account",
		Fields: []Field{
			{Name: ToCamelCase("phone"), Type: "int", Nullable: false, GormTag: "comment:'phone'", Validate: "required", JsonTag: "phone,omitempty"},
			{Name: ToCamelCase("name"), Type: "varchar", Nullable: false, GormTag: "comment:'用户名'", Validate: "required", JsonTag: "name,omitempty"},
			{Name: ToCamelCase("email"), Type: "varchar", Nullable: true, GormTag: "comment:'邮箱'", Validate: "email", JsonTag: "email,omitempty"},
			{Name: ToCamelCase("password"), Type: "varchar", Nullable: false, GormTag: "comment:'hash密码'", JsonTag: "password"},
			{Name: ToCamelCase("salt"), Type: "varchar", Nullable: false, GormTag: "comment:'加密盐'", JsonTag: "salt"},
		},
	}
	err := GenerateCode.Generate(table)
	if err != nil {
		t.Fatalf("生成文件错误%s", err)
	}
```
