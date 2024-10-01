package internal

import (
	// "context"
	// "fmt"
	"net"
)

func ServeGrpc(listener net.Listener, errch chan<- error) {
	_SLogger.Info("grpc server is up")

	var e error

	e = _RPCServer.Run(listener)
	errch <- e
}
