package ws

import (
	"bytes"
	"encoding/json"
	"fmt"
	// "log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Client struct {
	Id        uuid.UUID `json:"id"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"createdAt"`
	ClosedAt  time.Time `json:"closedAt"`
	PingAt    time.Time `json:"pingAt"`

	logger *zap.Logger
	conn   *websocket.Conn
	quit   chan struct{}
	once   *sync.Once
}

func (self Client) String() string {
	bts, _ := json.Marshal(self)
	return string(bts)
}

func (self *Client) Handle() (err error) {
loop:
	for {
		select {
		case <-self.quit:
			break loop
		default:
			// this can block the loop
			if err = self.HandleMessage(); err == nil {
				continue
			}

			switch err.(type) {
			// close 1006 (abnormal closure): unexpected EOF
			case *websocket.CloseError:
				err = nil
			default:
				self.logger.Error("handle_message", zap.Any("error", &err))
			}

			break loop
		}
	}

	self.Close()

	return err
}

func (self *Client) Close() {
	self.once.Do(func() {
		self.ClosedAt = time.Now()
		self.logger.Warn(
			"self_close",
			zap.String("pingAt", self.PingAt.Format(time.RFC3339)),
		)
		self.conn.Close()
	})
}

func (self *Client) HandleMessage() (err error) {
	var (
		ok        bool
		mt        int
		bts       []byte
		typ       string
		data, res map[string]any
	)

	defer func() {
		if data := recover(); data != nil {
			err = fmt.Errorf("read_error")
		}
	}()

	// var addr: net.Addr = conn.RemoteAddr()
	if mt, bts, err = self.conn.ReadMessage(); err != nil {
		return
	}

	defer func() {
		bts, _ = json.Marshal(res)
		err = self.conn.WriteMessage(mt, bts)
	}()

	data = make(map[string]any)

	if json.Unmarshal(bts, &data); err != nil {
		res = map[string]any{"type": "error", "msg": "unmarshal message error"}
		return nil
	}
	self.logger.Debug("received_message", zap.ByteString("msg", bytes.TrimSpace(bts)))

	if typ, ok = data["type"].(string); !ok {
		res = map[string]any{"type": "error", "messmsgage": "invalid field type"}
		return nil
	}

	switch typ {
	case "hello":
		name, _ := data["name"].(string)
		res = map[string]any{
			"type": "id", "id": self.Id, "msg": fmt.Sprintf("Welcome %s!", name),
		}
	case "id":
		res = map[string]any{"type": "id", "msg": "hello", "id": self.Id}
	default:
		res = data
	}

	return
}
