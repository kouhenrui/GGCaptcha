package GGCaptcha

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/kouhenrui/GGCaptcha/img"
	"github.com/kouhenrui/GGCaptcha/inter"
	"github.com/kouhenrui/GGCaptcha/store"
	"github.com/kouhenrui/GGCaptcha/utils"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type GGCaptcha struct {
	driver     inter.Driver           //验证码驱动
	store      inter.Store            //存储机制
	exptime    time.Duration          //过期时间
	windowTime time.Duration          //限制生成验证码的时间窗口
	limit      int                    //时间窗口内最大验证码数量
	mu         sync.Mutex             //并发安全控制
	requests   map[string][]time.Time //记录客户端请求时间戳
}

func NewGGCaptcha(driver inter.Driver, store inter.Store, t time.Duration, windowTime time.Duration, limit int) *GGCaptcha {
	return &GGCaptcha{
		driver:     driver,
		store:      store,
		exptime:    t,
		windowTime: windowTime,
		limit:      limit,
		requests:   make(map[string][]time.Time),
	}
}

func (g *GGCaptcha) checkLimit(clientID string, now time.Time) bool {
	timestamps := g.requests[clientID]

	//清理过期的请求记录
	newTimestamps := make([]time.Time, 0)
	for _, y := range timestamps {
		if now.Sub(y) <= g.windowTime {
			newTimestamps = append(newTimestamps, y)
		}
	}

	//更新记录
	g.requests[clientID] = newTimestamps

	//判断是否超过限制
	if len(newTimestamps) >= g.limit {
		return false
	}

	//添加当前请求时间戳
	g.requests[clientID] = append(g.requests[clientID], now)
	return true
}

func (g *GGCaptcha) GenerateGGCaptcha() (id, content string, err error) {
	content, answer, err := g.driver.GenerateDriverString()
	if err != nil {
		return "", "", err
	}
	id = utils.RandStr(8)
	err = g.store.Set(id, answer, g.exptime)
	if err != nil {
		return "", "", err
	}
	return id, content, nil
}

func (g *GGCaptcha) VerifyGGCaptcha(id, answer string, clear bool) bool {
	return g.store.Verify(id, answer, clear)
}

/*
 * @Title
 * @Description 使用算术验证
 * @Param
 * @return
 * @Author Acer
 * @Date 2024/9/14
 */

func (g *GGCaptcha) GenerateDriverMath() (id, content string, err error) {
	id = utils.RandStr(8)
	content, answer := g.driver.GenerateDriverMath()
	err = g.store.Set(id, answer, g.exptime)
	if err != nil {
		return "", "", err
	}
	return id, content, nil
}

/*
 * @Title
 * @Description 生成计算式图片验证码
 * @Param
 * @return
 * @Author Acer
 * @Date 2024/9/14
 */

func (g *GGCaptcha) GenerateDriverMathString() (id, content string, err error) {
	id = utils.RandStr(8)
	content, answer, err := g.driver.GenerateDriverMathString()
	if err != nil {
		return "", "", err
	}
	err = g.store.Set(id, answer, g.exptime)
	if err != nil {
		return "", "", err
	}
	return id, content, nil
}

/*
 * @Title
 * @Description 滑动验证码
 * @Param
 * @return
 * @Author Acer
 * @Date 2024/9/14
 */

func (g *GGCaptcha) GenerateDriverPuzzle() (id, bgImage, puzzleImage string, err error) {
	id = utils.RandStr(8)
	bgImage, puzzleImage, answer, err := g.driver.GenerateDriverPuzzle()
	if err != nil {
		return "", "", "", err
	}
	err = g.store.Set(id, strconv.Itoa(answer), g.exptime)
	if err != nil {
		return "", "", "", err
	}
	return id, bgImage, puzzleImage, nil
}

/*
 * @Title
 * @Description 刷新接口
 * @Param
 * @return
 * @Author Acer
 * @Date 2024/9/14
 */

func (g *GGCaptcha) RefreshCaptcha(ctype inter.CaptchaType) (id, content string, err error) {
	switch ctype {
	case inter.StringCaptcha:
		id, content, err = g.GenerateGGCaptcha()
	case inter.MathCaptcha:
		id, content, err = g.GenerateDriverMath()
	case inter.MathStringCaptcha:
		id, content, err = g.GenerateDriverMathString()
	}
	return id, content, err
}

type ImgOptions struct {
	Height       int         //图片高度
	Width        int         //图片宽度
	NoiseCount   int         //干扰线
	Count        int         //验证码数量
	Source       string      //数据源
	SourceLength int         //数据源长度
	SizePoint    float64     //字体大小
	Bgcolor      color.NRGBA //背景颜色
	FontStyle    font.Face   //字体
	UploadImg    image.Image

	MathNum      int //算术验证码数量
	PuzzleX      int // 滑动拼图的X位置
	PuzzleY      int // 滑动拼图的Y位置
	PuzzleWidth  int // 拼图块的宽度
	PuzzleHeight int // 拼图块的高度
}

func defaultImg() img.Img {
	var height = 80                                                           // 高度设置为80像素，确保有足够的空间容纳验证码和干扰
	var width = 240                                                           // 宽度设置为240像素，适合4-6个字符的验证码
	var noiseCount = 5                                                        // 5条干扰线，增强防破解性，但不影响阅读
	var count = 5                                                             // 验证码字符数量设置为5个，平衡安全性和用户友好性
	var source = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz0123456789" // 排除容易混淆的字符如 'O', '0', 'I', 'l'
	var sourceLength = len(source)                                            // 数据源字符长度
	var sizePoint = float64(height) * 0.6                                     // 字体大小设置为36，保证字符清晰度
	var fontColor = utils.RandColorRGBA(255)
	var bgColor = utils.RandColorRGBA(255)
	var fontStyle = utils.LoadDefaultFontFace()
	var mathNum = 3                                 // 默认生成两个数字进行加减法
	var puzzleWidth = width / 5                     // 拼图块占图片宽度的1/5
	var puzzleHeight = height / 5                   // 拼图块占图片高度的1/5
	var puzzleX = rand.Intn(width - width/5 - 10)   // 随机生成X坐标，预留边距
	var puzzleY = rand.Intn(height - height/5 - 10) // 随机生成Y坐标，预留边距
	return img.Img{
		Height:       height,             // 高度设置为80像素，确保有足够的空间容纳验证码和干扰
		Width:        width,              // 宽度设置为240像素，适合4-6个字符的验证码
		NoiseCount:   noiseCount,         // 5条干扰线，增强防破解性，但不影响阅读
		Count:        count,              // 验证码字符数量设置为5个，平衡安全性和用户友好性
		Source:       source,             // 排除容易混淆的字符如 'O', '0', 'I', 'l'
		SourceLength: sourceLength,       // 数据源字符长度
		SizePoint:    float64(sizePoint), // 字体大小设置为36，保证字符清晰度
		FontColor:    fontColor,
		BgColor:      bgColor,
		FontStyle:    fontStyle,
		UploadImg:    nil,
		MathNum:      mathNum,      //算术验证码数量
		PuzzleX:      puzzleWidth,  // 滑动拼图的X位置
		PuzzleY:      puzzleHeight, // 滑动拼图的Y位置
		PuzzleWidth:  puzzleX,      // 拼图块的宽度
		PuzzleHeight: puzzleY,      // 拼图块的高度

	}
}

/*
 * @Title
 * @Description load local img as background picture
 * @Param
 * @return
 * @Author Acer
 * @Date 2024/9/13
 */

func LoadLocalImg(imgPath string) img.Img {
	// 打开图片文件
	file, err := os.Open(imgPath)
	if err != nil {
		log.Fatalf("无法打开图片文件: %v", err)
	}
	defer file.Close()
	var imgs image.Image
	// 获取文件扩展名，判断图片格式
	ext := strings.ToLower(filepath.Ext(imgPath))
	switch ext {
	case ".jpg", ".jpeg":
		imgs, err = jpeg.Decode(file)
		if err != nil {
			log.Fatalf("JPEG 解码失败: %v", err)
		}
	case ".png":
		imgs, err = png.Decode(file)
		if err != nil {
			log.Fatalf("PNG 解码失败: %v", err)
		}
	default:
		log.Fatalf("不支持的图片格式: %s", ext)
	}
	return img.Img{UploadImg: imgs}
}

/*
 * @Title
 * @Description 创建绘图示例
 * @Param
 * @return
 * @Author Acer
 * @Date 2024/9/18
 */

func NewDriverString(imgOptions ...img.Img) *img.Img {
	var i img.Img
	if len(imgOptions) < 1 {
		i = defaultImg()
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

		if i.MathNum == 0 {
			i.MathNum = 3
		}
		if i.PuzzleWidth == 0 {
			i.PuzzleWidth = i.Width / 5
		}
		if i.PuzzleHeight == 0 {
			i.PuzzleHeight = i.Height / 5
		}
		if i.PuzzleX == 0 {
			i.PuzzleX = rand.Intn(i.Width - i.Width/5 - 10)
		}
		if i.PuzzleY == 0 {
			i.PuzzleY = rand.Intn(i.Height - i.Height/5 - 10)
		}
	}

	return &i
}

// Redis连接参数

type RedisOptions struct {
	Host     string
	Port     string
	Db       int
	PoolSize int
	MaxRetry int
}

/*
 * @Title
 * @Description 使用Redis作为存储
 * @Param
 * @return
 * @Author Acer
 * @Date 2024/9/18
 */

func NewRediStore(option RedisOptions) *store.RediStore {
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

/*
 * @Title
 * @Description 使用本地存储
 * @Param
 * @return
 * @Author Acer
 * @Date 2024/9/18
 */

func NewLocalStore() *store.LocalStore {
	return &store.LocalStore{
		Item: sync.Map{},
		Lock: sync.Mutex{},
	}
}
