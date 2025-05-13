package settings

import (
	// "fmt"

	"backend-api/internal/ws"

	"github.com/d2jvkpn/gotk"
	"github.com/d2jvkpn/gotk/ginx"
)

var (
	Logger   *gotk.ZapLogger
	WsServer *ws.Server
	JwtHMAC  *ginx.JwtHMAC
)
