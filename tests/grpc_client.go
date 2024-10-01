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

		conn   *grpc.ClientConn
		client *GrpcClient
		input  *proto.LogData
		res    *proto.LogId
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

	input = &proto.LogData{
		AppName:    "go-backend.tests",
		AppVersion: "0.1.0",

		RequestId:  uuid.New().String(),
		RequestAt:  time.Now().Format(gotk.RFC3339Milli),
		StatusCode: 200,

		LatencyMilli: 42,
		Identity:     map[string]string{"account": "test"},
		Data:         []byte(`{"module":"biz_user"}`),
	}
	if res, err = client.PushLog(_TestCtx, input); err != nil {
		return
	}

	log.Printf("<== response: %#v\n", res)
}
