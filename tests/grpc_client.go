package tests

import (
	"context"
	"flag"
	// "fmt"
	"log"

	"github.com/d2jvkpn/go-backend/proto"

	"github.com/d2jvkpn/gotk/cloud"
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
		ctx     context.Context

		conn   *grpc.ClientConn
		client *GrpcClient
		in     *proto.LogData
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

	ctx = context.TODO()

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

	in = &proto.LogData{
		AppName:    "go-backend/tests",
		AppVersion: "0.1.0",
		RequestId:  "testGrpcClient-01",
	}
	if res, err = client.PushLog(ctx, in); err != nil {
		return
	}

	log.Printf("==> response: %+#v\n", res)
}
