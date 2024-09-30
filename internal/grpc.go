package internal

import (
	// "context"
	// "fmt"
	"net"

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

func NewRPCServer(config *viper.Viper) (server *RPCServer, err error) {
	var (
		options []grpc.ServerOption
		uIntes  []grpc.UnaryServerInterceptor
		sIntes  []grpc.StreamServerInterceptor
	)

	options = make([]grpc.ServerOption, 0)
	uIntes = make([]grpc.UnaryServerInterceptor, 0)
	sIntes = make([]grpc.StreamServerInterceptor, 0)

	//
	if config.GetBool("otel") {
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
			config.GetString("cert"),
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

func SetupGrpc(config *viper.Viper) (err error) {
	var grpcConfig *viper.Viper

	grpcConfig = config.Sub("grpc")

	if _RPCServer, err = NewRPCServer(grpcConfig); err != nil {
		return err
	}

	grpc_health_v1.RegisterHealthServer(_RPCServer.Server, health.NewServer())

	// pkgXX.RegisterLogServiceServer(_RPCServer.Server, _RPCServer)

	return nil
}

func ServeGrpc(listener net.Listener, errch chan<- error) {
	_SLogger.Info("grpc server is up")

	var e error

	e = _RPCServer.Serve(listener)
	errch <- e
}
