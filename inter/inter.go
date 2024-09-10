package inter

import (
	"io"
	"time"
)

type GGCaptcha interface {
	GenerateGGCaptcha() (id, content string, err error)
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
}

// Item is captcha item inter
type Image interface {
	//WriteTo writes to a writer
	WriteTo(w io.Writer) (n int64, err error)
	//EncodeB64string encodes as base64 string
	EncodeB64string() string
}
