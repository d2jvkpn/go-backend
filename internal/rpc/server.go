package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net"

	"github.com/d2jvkpn/go-backend/proto"

	"github.com/google/uuid"
	grpcMdlw "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

type RPCServer struct {
	*grpc.Server
}

func NewRPCServer(config *viper.Viper) (server *RPCServer, err error) {
	var (
		serverOpts []grpc.ServerOption
		otelOpts   []otelgrpc.Option
		uIntes     []grpc.UnaryServerInterceptor
		sIntes     []grpc.StreamServerInterceptor
	)

	serverOpts = make([]grpc.ServerOption, 0)
	otelOpts = make([]otelgrpc.Option, 0)
	uIntes = make([]grpc.UnaryServerInterceptor, 0)
	sIntes = make([]grpc.StreamServerInterceptor, 0)

	// TODO: uIntes, sIntes
	if len(uIntes) > 0 {
		serverOpts = append(serverOpts,
			grpc.UnaryInterceptor(grpcMdlw.ChainUnaryServer(uIntes...)),
		)
	}

	if len(sIntes) > 0 {
		serverOpts = append(serverOpts,
			grpc.StreamInterceptor(grpcMdlw.ChainStreamServer(sIntes...)),
		)
	}

	if config.GetBool("trace") {
		otelOpts = append(otelOpts, otelgrpc.WithTracerProvider(otel.GetTracerProvider()))
	}

	if config.GetBool("meter") {
		otelOpts = append(otelOpts, otelgrpc.WithMeterProvider(otel.GetMeterProvider()))
	}

	if len(otelOpts) > 0 {
		serverOpts = append(
			serverOpts,
			grpc.StatsHandler(otelgrpc.NewServerHandler(otelOpts...)),
		)
	}

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
		serverOpts = append(serverOpts, grpc.Creds(creds))
	}

	server = new(RPCServer)
	server.Server = grpc.NewServer(serverOpts...)

	return server, nil
}

func (self *RPCServer) Run(listener net.Listener) (err error) {
	grpc_health_v1.RegisterHealthServer(self.Server, health.NewServer())

	proto.RegisterLogServiceServer(self.Server, self)

	return self.Serve(listener)
}

// biz
func (self *RPCServer) Push(ctx context.Context, record *proto.LogRequest) (
	*proto.LogResponse, error) {
	// TODO: biz

	fmt.Printf("<== PushLog: %+v\n", record)

	if !json.Valid(record.Data) {
		return nil, status.Error(codes.InvalidArgument, "invalid data: json bytes")
	}

	return &proto.LogResponse{Id: uuid.New().String(), EventId: record.GetEventId()}, nil
}
