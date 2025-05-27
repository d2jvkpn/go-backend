package utils

import (
	// "fmt"
	"math/rand"
	"time"

	"github.com/mojocn/base64Captcha"
	"github.com/spf13/viper"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Captcha struct {
	Enabled bool `mapstructure:"enabled"` // true
	Clear   bool `mapstructure:"clear"`   // true

	Height     int    `mapstructure:"height"`      // 96
	Width      int    `mapstructure:"width"`       // 160
	NoiseCount int    `mapstructure:"noise_count"` // 2
	ShowLine   int    `mapstructure:"show_line"`   // 2
	Length     int    `mapstructure:"length"`      // 5
	Source     string `mapstructure:"source"`      // ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789

	CollectNum int           `mapstructure:"collect_num"` // 100
	Duration   time.Duration `mapstructure:"duration"`    // 5m

	*base64Captcha.Captcha `mapstructure:"-"`
}

func CaptchaFromViper(vp *viper.Viper) (captcha *Captcha, err error) {
	captcha = new(Captcha)

	if err = vp.Unmarshal(captcha); err != nil {
		return nil, err
	}

	if err = captcha.Setup(); err != nil {
		return nil, err
	}

	// fmt.Printf("~~~ %t\n",  captcha.Clear)

	return captcha, nil
}

func (self *Captcha) Setup() (err error) {
	var (
		id     string
		answer string
		driver base64Captcha.Driver
		store  base64Captcha.Store
	)

	/*
		base64Captcha.DefaultEmbeddedFonts.LoadFontByName()

		fonts := []string{
			"3Dumb.ttf", "ApothecaryFont.ttf", "Comismsh.ttf", "DENNEthree-dee.ttf", "DeborahFancyDress.ttf",
			"Flim-Flam.ttf", "RitaSmith.ttf", "actionj.ttf", "chromohv.ttf",
		}
	*/

	driver = base64Captcha.NewDriverString(
		// 60, 100, 200, 2, 6,
		self.Height,
		self.Width,
		self.NoiseCount,
		// base64Captcha.OptionShowHollowLine,
		self.ShowLine,
		self.Length,
		self.Source,
		nil,
		nil,
		nil,
	)
	store = base64Captcha.NewMemoryStore(self.CollectNum, self.Duration)

	self.Captcha = base64Captcha.NewCaptcha(driver, store)

	if id, _, answer, err = self.Captcha.Generate(); err != nil {
		return err
	}
	_ = self.Captcha.Verify(id, answer, true)

	return nil
}

func (self *Captcha) Verify(id, answer string, clear bool) bool {
	if !self.Enabled {
		return true
	}

	if id == "" || answer == "" {
		return false
	}

	return self.Captcha.Verify(id, answer, clear)
}
