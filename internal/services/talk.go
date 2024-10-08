package services

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/d2jvkpn/go-backend/internal/settings"
	"github.com/d2jvkpn/go-backend/internal/ws"
)

func talk(ctx *gin.Context) {
	var (
		err    error
		conn   *websocket.Conn
		client *ws.Client
	)

	defer func() {
		if err != nil {
			log.Printf("!!! error: %v\n", err)
		}
	}()

	conn, err = settings.WsServer.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}
	// defer conn.Close()

	client = settings.WsServer.NewClient(ctx, conn)

	// to avoid dead lock when for loop blocked by HandleMessage, don't use an unbuffered channel
	err = client.Handle()
}
