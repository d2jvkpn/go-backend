package settings

import (
	// "fmt"

	"github.com/d2jvkpn/go-backend/internal/ws"

	"github.com/d2jvkpn/gotk"
)

var (
	Logger   *gotk.ZapLogger
	WsServer *ws.Server
)
