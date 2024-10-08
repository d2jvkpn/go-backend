package tests

import (
	"bufio"
	"bytes"
	"flag"
	// "fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

func testWSClient(args []string) {
	var (
		addr    string
		err     error
		flagSet *flag.FlagSet
		client  *WsClient
	)

	flagSet = flag.NewFlagSet("ws_client", flag.ExitOnError)
	flagSet.StringVar(&addr, "addr", "ws://127.0.0.1:9011/socket/talk", "websocket address")
	flagSet.Parse(args)

	if client, err = NewWsClient(addr); err != nil {
		log.Fatal(err)
	}

	client.SetPing(5 * time.Second)
	if err = client.ReadMsgFromPipeFile("data/ws_client.fifo"); err != nil {
		client.Close()
		log.Fatal(err)
	}

	client.HandleMessage()
	client.Close()
}

type WsClient struct {
	conn  *websocket.Conn
	mutex *sync.Mutex
	done  chan struct{}
	once  *sync.Once
}

func NewWsClient(addr string) (client *WsClient, err error) {
	var conn *websocket.Conn

	if !strings.HasPrefix(addr, "ws") {
		addr = "ws://" + addr
	}

	if conn, _, err = websocket.DefaultDialer.Dial(addr, nil); err != nil {
		return nil, err
	}

	client = &WsClient{
		conn:  conn,
		mutex: new(sync.Mutex),
		done:  make(chan struct{}, 1),
		once:  new(sync.Once),
	}

	// overwrite default handler(when receive a ping)
	//conn.SetPingHandler(func(data string) (err error) {
	//	log.Printf("<~~ Recv ping: %q\n", data)

	//	client.mutex.Lock()
	//	err = conn.WriteMessage(websocket.PongMessage, nil)
	//	client.mutex.Unlock()

	//	if err != nil {
	//		log.Printf("!!! WriteMessage Pong: %v\n", err)
	//	}

	//	return err
	//})

	// overwrite default handler(when receive a pong after send a ping)
	//conn.SetPongHandler(func(data string) (err error) {
	//	log.Printf("<~~ recv pong: %q\n", data)
	//	return nil
	//})

	return client, nil
}

func (self *WsClient) Close() {
	self.once.Do(func() {
		close(self.done)
		self.conn.Close()
	})
}

func (self *WsClient) SetPing(dur time.Duration) {
	ticker := time.NewTicker(dur)

	go func() {
		var pingId uint64

	loop:
		for {
			select {
			case <-self.done:
				break loop
			case _ = <-ticker.C:
				pingId += 1
				data := []byte(strconv.FormatUint(pingId, 10))
				self.mutex.Lock()

				log.Printf("~~> send ping: %q\n", data)
				err := self.conn.WriteMessage(websocket.PingMessage, []byte(data))
				if err != nil {
					log.Printf("!!! send ping: %v\n", err)
				}
				self.mutex.Unlock()
			}
		}
	}()
}

func (self *WsClient) ReadMsgFromPipeFile(fp string) (err error) {
	var (
		file   *os.File
		reader *bufio.Reader
	)

	_ = os.Remove(fp)
	if err = os.MkdirAll(filepath.Dir(fp), 0755); err != nil {
		return err
	}
	if err = syscall.Mkfifo(fp, 0640); err != nil {
		return err
	}

	if file, err = os.OpenFile(fp, os.O_CREATE, os.ModeNamedPipe); err != nil {
		return err
	}

	reader = bufio.NewReader(file)

	go func() {
		var (
			bts []byte
			msg string
		)

		log.Println("~~~ Read from pipefile:", fp)

		for {
			select {
			case <-self.done:
				return
			default:
				// pass
			}

			if msg, err = reader.ReadString('\n'); err != nil {
				if err == io.EOF {
					continue
				}

				log.Printf("!!! ReadString: %v\n", err)
				break
			}

			msg = strings.TrimSpace(msg)

			if msg == "\\q" {
				log.Println("!!! exit client")
				self.Close()
				break
			}

			if bts = []byte(msg); len(bts) == 0 {
				continue
			}

			self.mutex.Lock()
			err = self.conn.WriteMessage(websocket.TextMessage, bts)
			self.mutex.Unlock()

			if err != nil {
				log.Printf("!!! WriteMessage text: %v\n", err)
				break
			}
		}

		_ = file.Close()
		_ = os.Remove(fp)
	}()

	return nil
}

func (self *WsClient) HandleMessage() {
	var (
		typ int
		bts []byte
		err error
	)

	for {
		if typ, bts, err = self.conn.ReadMessage(); err != nil {
			log.Printf("!!! ReadMessage error: %[1]T, %[1]v\n", err)
			break
		}

		log.Printf("<-- ReadMessage: type=%d, %s\n", typ, bytes.TrimSpace(bts))
		if typ == websocket.CloseMessage {
			break
		}
	}

	self.Close()
}

// https://pkg.go.dev/github.com/gorilla/websocket#pkg-types
//	TextMessage = 1
//	BinaryMessage = 2
//	CloseMessage = 8
//	PingMessage = 9
//	PongMessage = 10
