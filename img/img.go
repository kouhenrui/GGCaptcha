package img

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/fogleman/gg"
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
	"strings"
)

type Img struct {
	Height       int    //图片高度
	Width        int    //图片宽度
	NoiseCount   int    //干扰线
	Count        int    //验证码数量
	Source       string //数据源
	SourceLength int    //数据源长度
	FontColor    color.Color
	SizePoint    float64     //字体大小
	BgColor      color.NRGBA //背景颜色
	FontStyle    font.Face   //字体
	UploadImg    image.Image

	MathNum      int //算术验证码数量
	PuzzleX      int // 滑动拼图的X位置
	PuzzleY      int // 滑动拼图的Y位置
	PuzzleWidth  int // 拼图块的宽度
	PuzzleHeight int // 拼图块的高度
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

func DefaultImg() Img {
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
	return Img{
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

func LoadLocalImg(imgPath string) Img {
	// 打开图片文件
	file, err := os.Open(imgPath)
	if err != nil {
		log.Fatalf("无法打开图片文件: %v", err)
	}
	defer file.Close()
	var img image.Image
	// 获取文件扩展名，判断图片格式
	ext := strings.ToLower(filepath.Ext(imgPath))
	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
		if err != nil {
			log.Fatalf("JPEG 解码失败: %v", err)
		}
	case ".png":
		img, err = png.Decode(file)
		if err != nil {
			log.Fatalf("PNG 解码失败: %v", err)
		}
	default:
		log.Fatalf("不支持的图片格式: %s", ext)
	}
	return Img{UploadImg: img}
}
func NewDriverString(imgOptions ...Img) *Img {
	var i Img
	if len(imgOptions) < 1 {
		i = DefaultImg()
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

/*
 * @Title
 * @Description 根据默认值或图片地址生成验证码
 * @Param
 * @return
 * @Author Acer
 * @Date 2024/9/13
 */

func (m *Img) GenerateDriverString() (content, answer string, err error) {
	//generate rand string
	textPlain := utils.RandStr(m.Count)

	var dc *gg.Context
	if m.UploadImg != nil {
		dc = gg.NewContextForImage(m.UploadImg)
	} else {
		dc = m.makeBGColor()
	}

	dc.SetFontFace(m.FontStyle) //设置字体格式
	dc.SetColor(m.FontColor)    //设置字体颜色
	m.writeText(dc, textPlain)  //绘制验证码图片
	m.interfereLine(dc)         //绘制干扰线
	//dc.Fill()                   //填充填充路径的内部区域 和stroke方法合用
	var buffer bytes.Buffer

	err = dc.EncodePNG(&buffer)
	if err != nil {
		return "", "", err
	}
	content = base64.StdEncoding.EncodeToString(buffer.Bytes())
	return content, textPlain, nil
}

/*
 * @Title
 * @Description
 * @Param
 * @return
 * @Author Acer
 * @Date 2024/9/13
 */

func (m *Img) GenerateDriverMathString() (content, answer string, err error) {
	expression, result := utils.RandMath(m.MathNum)
	var dc *gg.Context
	if m.UploadImg != nil {
		dc = gg.NewContextForImage(m.UploadImg)
	} else {
		dc = m.makeBGColor()
	}

	dc.SetFontFace(m.FontStyle) //设置字体格式
	dc.SetColor(m.FontColor)    //设置字体颜色
	m.writeText(dc, expression) //绘制验证码图片
	m.interfereLine(dc)         //绘制干扰线
	//dc.Fill()                   //填充填充路径的内部区域 和stroke方法合用
	var buffer bytes.Buffer

	err = dc.EncodePNG(&buffer)
	if err != nil {
		return "", "", err
	}
	content = base64.StdEncoding.EncodeToString(buffer.Bytes())
	return content, result, nil
}

/*
 * @Title
 * @Description 生成纯运算符式验证码
 * @Param
 * @return
 * @Author Acer
 * @Date 2024/9/13
 */

func (m *Img) GenerateDriverMath() (content, answer string) {
	expression, result := utils.RandMath(m.MathNum)
	return expression, result
}

/*
 * @Title
 * @Description 生成滑动式验证码
 * @Param
 * @return
 * @Author Acer
 * @Date 2024/9/13
 */

func (m *Img) GenerateDriverPuzzle() (bgImage string, puzzleImage string, targetX int, err error) {
	var dc *gg.Context
	var puzzlePiece image.Image
	if m.UploadImg != nil {
		dc = gg.NewContextForImage(m.UploadImg)
		puzzleRect := image.Rect(m.PuzzleX, m.PuzzleY, m.PuzzleX+m.PuzzleWidth, m.PuzzleY+m.PuzzleHeight)
		puzzlePiece = m.UploadImg.(interface {
			SubImage(r image.Rectangle) image.Image
		}).SubImage(puzzleRect)

	} else {
		dc = m.puzzleString()
		puzzleRect := image.Rect(m.PuzzleX, m.PuzzleY, m.PuzzleX+m.PuzzleWidth, m.PuzzleY+m.PuzzleHeight)
		// 使用背景图生成拼图块
		puzzlePiece = dc.Image().(interface {
			SubImage(r image.Rectangle) image.Image
		}).SubImage(puzzleRect)
	}
	// 绘制拼图块
	dcPiece := gg.NewContextForImage(puzzlePiece)
	var puzzleBuffer bytes.Buffer
	err = dcPiece.EncodePNG(&puzzleBuffer)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to encode puzzle piece: %w", err)
	}
	puzzleImage = base64.StdEncoding.EncodeToString(puzzleBuffer.Bytes())

	// 在背景图上挖出拼图块
	dc.SetRGBA(1, 1, 1, 1) // 设置填充为白色
	dc.DrawRectangle(float64(m.PuzzleX), float64(m.PuzzleY), float64(m.PuzzleWidth), float64(m.PuzzleHeight))
	dc.Fill()
	// 转换背景图为 base64
	var bgBuffer bytes.Buffer
	err = dc.EncodePNG(&bgBuffer)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to encode background image: %w", err)
	}
	bgImage = base64.StdEncoding.EncodeToString(bgBuffer.Bytes())

	// 随机生成目标 X 位置
	targetX = m.PuzzleX + rand.Intn(m.PuzzleWidth)

	return bgImage, puzzleImage, targetX, nil
}

/*
 * @Title
 * @Description 生成模拟背景
 * @Param
 * @return
 * @Author Acer
 * @Date 2024/9/13
 */
func (m *Img) puzzleString() (dc *gg.Context) {
	textPlain := utils.RandStr(m.Count)
	dc = m.makeBGColor()
	dc.SetFontFace(m.FontStyle) //设置字体格式
	dc.SetColor(m.FontColor)    //设置字体颜色
	m.writeText(dc, textPlain)  //绘制验证码图片
	m.interfereLine(dc)         //绘制干扰线

	return dc
}

/*
 * @Title
 * @Description 设置图片干扰线
 * @Param
 * @return
 * @Author Acer
 * @Date 2024/9/10
 */
func (m *Img) interfereLine(dc *gg.Context) {
	// 干扰线
	if m.NoiseCount > 0 {
		if m.NoiseCount > 8 {
			m.NoiseCount = 8
		}
		for i := 0; i < m.NoiseCount; i++ {
			x1, y1 := utils.RandPos(m.Width, m.Height)
			x2, y2 := utils.RandPos(m.Width, m.Height)
			r, g, b, a := utils.RandColor(255)
			w := float64(rand.Intn(3) + 1)
			dc.SetRGBA255(r, g, b, a)
			dc.SetLineWidth(w)
			dc.DrawLine(x1, y1, x2, y2)
			dc.Stroke() //绘制轮廓
		}
	}

}

/*
 * @Title
 * @Description 渲染文字到图片上
 * @Param
 * @return
 * @Author Acer
 * @Date 2024/9/10
 */
func (m *Img) writeText(dc *gg.Context, text string) {
	log.Println(text)
	// 获取文字的字符数组
	characters := []rune(text)
	charCount := len(characters)
	log.Println(charCount, "字符长度")
	// 每个字符的水平间距
	charSpacing := float64(m.Width) / float64(charCount+1)
	for i := 0; i < charCount; i++ {
		//r, g, b, _ := utils.RandColor(100)
		//dc.SetRGBA255(r, g, b, 255)

		// 字符的x坐标，均匀分布
		x := charSpacing * float64(i+1)

		// 字符的y坐标，居中但带有轻微随机浮动
		y := float64(m.Height) * 0.5
		yfload := 5 - rand.Float64()*10 + y // 上下浮动5个像素

		// 为字符增加随机旋转角度，增强干扰效果
		rotation := (rand.Float64() * 30) - 15 // 每个字符随机旋转 -15 到 15 度
		dc.RotateAbout(gg.Radians(rotation), x, y)

		// 绘制字符到图片上，居中显示
		dc.DrawStringAnchored(string(characters[i]), x, yfload, 0.5, 0.5)

		// 恢复原来的旋转角度，避免影响后面的字符
		dc.RotateAbout(gg.Radians(-rotation), x, y)
		dc.Stroke()

	}
	//dc.Clear()
}

/*
 * @Title
 * @Description 绘制背景板
 * @Param
 * @return
 * @Author Acer
 * @Date 2024/9/13
 */
func (m *Img) makeBGColor() (dc *gg.Context) {
	dc = gg.NewContext(m.Width, m.Height) //设置图片大小
	bgR, bgG, bgB, bgA := utils.RandColorRange(5, 255)
	dc.SetRGBA255(bgR, bgG, bgB, bgA) //设置初始背景板颜色
	dc.Clear()                        //将背景板填充为之前设置颜色
	return dc
}
