package img

import (
	"GGCaptcha/utils"
	"bytes"
	"github.com/fogleman/gg"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"image/color"
	"math/rand"
)

type Img struct {
	Height       int         //图片高度
	Width        int         //图片宽度
	NoiseCount   int         //干扰线
	Count        int         //验证码数量
	Source       string      //数据源
	SourceLength int         //s=数据源长度
	SizePoint    float64     //字体大小
	bgcolor      color.NRGBA //背景颜色
	FontStyle    font.Face   //字体
}

func NewDriverString(imgOptions ...Img) *Img {
	var img Img
	if len(imgOptions) < 1 {
		img.Height = 60
		img.Width = 120
		img.Count = 4
		img.SizePoint = 20
		img.NoiseCount = 4
		img.Source = "ABCDEFGHIJKMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz0123456789"
		img.SourceLength = len(img.Source)
		img.bgcolor = utils.RandColorRGBA(255)
		img.FontStyle = loadFontFace(img.SizePoint)
	} else {
		img = imgOptions[0]
	}

	return &img
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

// 加载字体
func loadFontFace(sizePoint float64) font.Face {
	// 使用 Go 自带的字体，也可以使用本地字体文件
	ttf, err := opentype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}
	face, err := opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    sizePoint,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	return face
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
