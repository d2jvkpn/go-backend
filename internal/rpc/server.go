package rpc

import (
	"context"
	"fmt"
	"net"

	"github.com/d2jvkpn/go-backend/proto"

	grpcMdlw "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type RPCServer struct {
	*grpc.Server
}

func NewRPCServer(config *viper.Viper, enableOtel bool) (server *RPCServer, err error) {
	var (
		options []grpc.ServerOption
		uIntes  []grpc.UnaryServerInterceptor
		sIntes  []grpc.StreamServerInterceptor
	)

	options = make([]grpc.ServerOption, 0)
	uIntes = make([]grpc.UnaryServerInterceptor, 0)
	sIntes = make([]grpc.StreamServerInterceptor, 0)

	if enableOtel {
		uIntes = append(
			uIntes,
			otelgrpc.UnaryServerInterceptor( /*opts ...Option*/ ),
		)

		sIntes = append(
			sIntes,
			otelgrpc.StreamServerInterceptor( /*opts ...Option*/ ),
		)
	}

	options = append(options,
		grpc.UnaryInterceptor(grpcMdlw.ChainUnaryServer(uIntes...)),
		grpc.StreamInterceptor(grpcMdlw.ChainStreamServer(sIntes...)),
	)

	//
	if config.GetBool("tls") {
		var creds credentials.TransportCredentials

		creds, err = credentials.NewServerTLSFromFile(
			config.GetString("cer"),
			config.GetString("key"),
		)
		if err != nil {
			return nil, err
		}
		options = append(options, grpc.Creds(creds))
	}

	server = new(RPCServer)
	server.Server = grpc.NewServer(options...)

	return server, nil
}

func (self *RPCServer) Run(listener net.Listener) (err error) {
	grpc_health_v1.RegisterHealthServer(self.Server, health.NewServer())

	proto.RegisterLogServiceServer(self.Server, self)

	return self.Serve(listener)
}

// biz
func (self *RPCServer) PushLog(ctx context.Context, record *proto.LogData) (*proto.LogId, error) {
	// TODO: biz

	fmt.Printf("<== PushLog: %+v\n", record)

	return &proto.LogId{Id: record.GetRequestId()}, nil
}
