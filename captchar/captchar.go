package captchar

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/color"
	"image/png"
	"strings"
	"time"

	"github.com/afocus/captcha"
	"github.com/dchest/uniuri"
	"github.com/empirefox/esecend/front"
	"github.com/patrickmn/go-cache"
)

var (
	expires = 2 * time.Minute
	clears  = 10 * time.Minute
)

type Captchar interface {
	New(userId uint) (*front.Captcha, error)
	Verify(userId uint, key, value string) bool
}

type captchar struct {
	cap      *captcha.Captcha
	capCache *cache.Cache
}

func NewCaptchar(font string) (Captchar, error) {
	cap := captcha.New()

	if err := cap.SetFont(font); err != nil {
		return nil, err
	}

	cap.SetSize(92, 32)
	cap.SetDisturbance(captcha.MEDIUM)
	cap.SetBkgColor(color.RGBA{255, 255, 255, 255})
	cap.SetFrontColor(
		color.RGBA{255, 153, 18, 255},
		color.RGBA{128, 128, 192, 255},
		color.RGBA{232, 208, 152, 255},
		color.RGBA{41, 36, 33, 255},
		color.RGBA{51, 102, 153, 255},
		color.RGBA{102, 153, 204, 255},
		color.RGBA{180, 91, 62, 255},
		color.RGBA{0, 178, 113, 255},
	)

	return &captchar{
		cap:      cap,
		capCache: cache.New(expires, clears),
	}, nil
}

func (c *captchar) New(userId uint) (*front.Captcha, error) {
	img, value := c.cap.Create(4, captcha.ALL)

	var b bytes.Buffer
	b.WriteByte('"')
	b64 := base64.NewEncoder(base64.StdEncoding, &b)
	defer b64.Close()
	if err := png.Encode(b64, img); err != nil {
		return nil, err
	}
	b.WriteByte('"')

	key := uniuri.New()
	for c.capCache.Add(key, fmt.Sprintln("%d:%s", userId, strings.ToLower(value)), cache.DefaultExpiration) != nil {
		key = uniuri.NewLen(20)
	}

	data := json.RawMessage(b.Bytes())
	return &front.Captcha{
		ID:     key,
		Base64: &data,
	}, nil
}

func (c *captchar) Verify(userId uint, key, value string) bool {
	fact, ok := c.capCache.Get(key)
	if ok {
		c.capCache.Delete(key)
	}
	return ok && fact.(string) == fmt.Sprintln("%d:%s", userId, strings.ToLower(value))
}
