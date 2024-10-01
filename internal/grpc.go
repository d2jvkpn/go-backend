package internal

import (
	// "context"
	// "fmt"
	"net"

	"go.uber.org/zap"
)

func ServeGrpc(listener net.Listener, errch chan<- error) {
	_SLogger.Debug("grpc server is up")

	var e error

	if e = _RPCServer.Run(listener); e != nil {
		_Logger.Error("grpc server has been shutdown", zap.String("error", e.Error()))
		errch <- e
	} else {
		_Logger.Info("grpc server has been shutdown")
		errch <- nil
	}
}
