package img

import (
	"GGCaptcha/utils"
	"bytes"
	"github.com/fogleman/gg"
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
}

func defaultImg() Img {
	return Img{
		Height:       60,
		Width:        120,
		Count:        4,
		SizePoint:    20,
		NoiseCount:   4,
		Source:       "ABCDEFGHIJKMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz0123456789",
		SourceLength: len("ABCDEFGHIJKMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz0123456789"),
		Bgcolor:      utils.RandColorRGBA(255),
		FontStyle:    utils.LoadDefaultFontFace(),
		UploadImg:    nil,
	}
}

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
		i = defaultImg()
	} else {
		// 使用传入的自定义配置
		i = imgOptions[0]

		//// 如果没有设置某些字段，则为其提供默认值
		//if i.Height == 0 {
		//	i.Height = 60
		//}
		//if i.Width == 0 {
		//	i.Width = 120
		//}
		//if i.Count == 0 {
		//	i.Count = 4
		//}
		//if i.SizePoint == 0 {
		//	i.SizePoint = 20
		//}
		//if i.NoiseCount == 0 {
		//	i.NoiseCount = 4
		//}
		//if i.Source == "" {
		//	i.Source = "ABCDEFGHIJKMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz0123456789"
		//	i.SourceLength = len(i.Source)
		//}
		//if i.Bgcolor == (color.NRGBA{}) {
		//	i.Bgcolor = utils.RandColorRGBA(255)
		//}
		//if i.FontStyle == nil {
		//	i.FontStyle = utils.LoadFontFace(i.SizePoint)
		//}
	}

	return &i
}
func (m *Img) GenerateDriverString() (content, answer string, err error) {
	dc := gg.NewContext(m.Width, m.Height)
	bgR, bgG, bgB, bgA := utils.RandColorRange(100, 255)
	dc.SetRGBA255(bgR, bgG, bgB, bgA)
	dc.Clear()
	dc.SetFontFace(m.FontStyle)
	textPlain := utils.RandStr(m.Count)
	m.writeText(dc, textPlain)
	dc.Fill()
	var buffer bytes.Buffer
	err = dc.EncodePNG(&buffer)
	if err != nil {
		return "", "", err
	}
	content = buffer.String()
	return content, textPlain, nil
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
			dc.Stroke()
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
	for i := 0; i < m.SourceLength; i++ {
		r, g, b, _ := utils.RandColor(100)
		dc.SetRGBA255(r, g, b, 255)
		x := float64(m.Width/m.SourceLength*i) + m.SizePoint*0.6
		y := float64(m.Height / 2)
		xfload := 5 - rand.Float64()*10 + x
		yfload := 5 - rand.Float64()*10 + y
		radians := 40 - rand.Float64()*80
		dc.RotateAbout(gg.Radians(radians), x, y)
		dc.DrawStringAnchored(text, xfload, yfload, 0.2, 0.5)
		dc.RotateAbout(-1*gg.Radians(radians), x, y)
		dc.Stroke()
	}

}

func GenerateDriverImage(im image.Image) {
	gg.NewContextForImage(im)

}
func GenerateDriverRGBA(im *image.RGBA) {
	gg.NewContextForRGBA(im)
}
