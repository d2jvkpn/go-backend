package settings

import (
	// "fmt"

	"backend-api/internal/ws"
	"backend-api/pkg/utils"

	"github.com/d2jvkpn/gotk"
	"github.com/d2jvkpn/gotk/ginx"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var (
	Config   *viper.Viper
	Project  *viper.Viper
	Logger   *gotk.ZapLogger
	Captcha  *utils.Captcha
	WsServer *ws.Server
	JwtHMAC  *ginx.JwtHMAC
	Redis    *redis.Client
)
