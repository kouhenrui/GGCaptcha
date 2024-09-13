package inter

import (
	"io"
	"time"
)

type GGCaptcha interface {
	GenerateGGCaptcha() (id, content, answer string, err error)
	VerifyGGCaptcha(id, answer string, clear bool) bool
}
type Store interface {
	Set(id string, value string, t time.Duration) error
	Exist(id string) bool
	Get(id string, clear bool) (string, error)
	Verify(id, answer string, clear bool) bool
}
type Driver interface {
	GenerateDriverString() (content, answer string, err error)
	GenerateDriverMathString() (content, answer string, err error)
	GenerateDriverMath() (content, answer string)
	GenerateDriverPuzzle() (bgImage string, puzzleImage string, targetX int, err error)
}

// Item is captcha item inter
type Image interface {
	//WriteTo writes to a writer
	WriteTo(w io.Writer) (n int64, err error)
	//EncodeB64string encodes as base64 string
	EncodeB64string() string
}
