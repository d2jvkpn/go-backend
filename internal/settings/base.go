package settings

import (
	// "fmt"

	"github.com/d2jvkpn/go-backend/internal/ws"

	"github.com/d2jvkpn/gotk"
	"github.com/d2jvkpn/gotk/ginx"
)

var (
	Logger   *gotk.ZapLogger
	WsServer *ws.Server
	JwtHMAC  *ginx.JwtHMAC
)
