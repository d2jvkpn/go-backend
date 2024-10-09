package ws

import (
	// "fmt"
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Server struct {
	clients   map[uuid.UUID]*Client
	broadcast chan []byte
	mutex     sync.Mutex
	Upgrader  *websocket.Upgrader
}

func NewServer() *Server {
	return &Server{
		clients:   make(map[uuid.UUID]*Client),
		broadcast: make(chan []byte),
		Upgrader: &websocket.Upgrader{
			EnableCompression: true,
			HandshakeTimeout:  2 * time.Second,
			ReadBufferSize:    1024,
			WriteBufferSize:   1024,
			// CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

func (self *Server) NewClient(ctx *gin.Context, conn *websocket.Conn) *Client {
	client := &Client{
		Id:        uuid.New(),
		Address:   ctx.ClientIP(),
		CreatedAt: time.Now(),

		conn: conn,
		quit: make(chan struct{}, 1),
		once: new(sync.Once),
	}

	conn.SetPingHandler(func(data string) (err error) {
		// log.Printf("~~~ %s ping: %q\n", client.Id, data)
		client.PingAt = time.Now()

		return conn.WriteMessage(websocket.PongMessage, []byte(data))
	})

	conn.SetCloseHandler(func(code int, text string) error {
		log.Printf("<== %s closed: code=%d, text=%q\n", client.Id, code, text)
		client.quit <- struct{}{}

		return nil
	})

	self.mutex.Lock()
	self.clients[client.Id] = client
	self.mutex.Unlock()

	return client
}

func (self *Server) RemoveClient(id uuid.UUID, mannual bool) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if mannual {
		if client, ok := self.clients[id]; ok {
			client.Close()
		}
	}

	delete(self.clients, id)
}

func (self *Server) Shutdown() {
	if self == nil {
		return
	}

	self.mutex.Lock()
	defer self.mutex.Unlock()

	for id, client := range self.clients {
		client.Close()
		delete(self.clients, id)
	}
}
