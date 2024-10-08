package tests

import (
	// "context"
	"flag"
	// "fmt"
	"log"
	"time"

	"github.com/d2jvkpn/go-backend/proto"

	"github.com/d2jvkpn/gotk"
	"github.com/d2jvkpn/gotk/cloud"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type GrpcClient struct {
	conn *grpc.ClientConn
	proto.LogServiceClient
}

func NewGrpcClient(conn *grpc.ClientConn) *GrpcClient {
	return &GrpcClient{
		conn:             conn,
		LogServiceClient: proto.NewLogServiceClient(conn),
	}
}

func testGrpcClient(args []string) {
	var (
		addr    string
		tls     bool
		err     error
		flagSet *flag.FlagSet

		conn     *grpc.ClientConn
		client   *GrpcClient
		request  *proto.LogRequest
		response *proto.LogResponse
	)

	flagSet = flag.NewFlagSet("testGrpcClient", flag.ContinueOnError)

	flagSet.StringVar(&addr, "addr", "localhost:9016", "grpc address")
	flagSet.BoolVar(&tls, "tls", false, "enable tls")
	flagSet.Parse(args)

	defer func() {
		if conn != nil {
			conn.Close()
		}

		if err != nil {
			log.Fatal(err)
		}
	}()

	inte := cloud.HeaderInterceptor{
		Headers: map[string]string{"hello": "world"},
	}

	opts := []grpc.DialOption{
		grpc.WithUnaryInterceptor(inte.Unary()),
		grpc.WithStreamInterceptor(inte.Stream()),
	}
	if !tls {
		opts = append(opts, grpc.WithInsecure())
	}

	if conn, err = grpc.Dial(addr, opts...); err != nil {
		log.Fatal(err)
	}
	log.Println("==> grpc connected:", addr)

	client = NewGrpcClient(conn)

	request = &proto.LogRequest{
		EventLevel: proto.EventLevel_info,
		AppName:    "go-backend/tests",
		AppVersion: "0.1.0",

		Service: "http",
		Id:      uuid.New().String(),
		At:      time.Now().Format(gotk.RFC3339Milli),
		BizName: "POST@/api/v1/open/login",
		BizData: map[string]string{
			"query":  "region=cn&city=shanghai",
			"status": "OK",
			"client": "web",
		},
		Identities: map[string]string{"account": "test", "role": "normal"},
		Code:       "ok",

		LatencyMilli: 42,
		Data:         []byte(`{"module":"biz_user"}`),
	}

	log.Printf("==> send: %#v\n", request)
	if response, err = client.Push(_TestCtx, request); err != nil {
		return
	}

	log.Printf("<== response: %#v\n", response)
}
