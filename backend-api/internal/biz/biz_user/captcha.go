package biz_user

import (
	"fmt"
	"time"

	"backend-api/internal/settings"
	"backend-api/pkg/structs"

	"github.com/d2jvkpn/errx"
)

type CaptchaResponse struct {
	// required: true
	Enabled bool   `json:"enabled" extensions:"x-order=01"`
	Id      string `json:"id,omitempty" extensions:"x-order=02"`

	Length      int    `json:"length,omitempty" extensions:"x-order=03"`
	Base64Image string `json:"base64Image,omitempty" extensions:"x-order=03"` // data:image/png;base64,...
	ExpiresAt   int64  `json:"expiresAt,omitempty" extensions:"x-order=04"`   // unix timstamp
}

func NewCaptcha() (item *CaptchaResponse, err *errx.ErrX) {
	var (
		e   error
		now time.Time
	)

	item = new(CaptchaResponse)
	if item.Enabled = settings.Captcha.Enabled; !item.Enabled {
		return item, nil
	}
	item.Length = settings.Captcha.Length

	now = time.Now()
	item.Id, item.Base64Image, _, e = settings.Captcha.Generate()
	if e != nil {
		return nil, structs.InternalError(e).WithCode("new_captcha")
	}
	item.ExpiresAt = now.Add(settings.Captcha.Duration).Unix()

	return item, nil
}

func VerifyCaptcha(captchaId, captchaAnswer string) (err *errx.ErrX) {
	if !settings.Captcha.Verify(captchaId, captchaAnswer, settings.Captcha.Clear) {
		const msg = "captcha verify failed"
		return structs.BizError(fmt.Errorf(msg)).WithCode("captcha_verify_failed").WithMsg(msg)
	}

	return nil
}
