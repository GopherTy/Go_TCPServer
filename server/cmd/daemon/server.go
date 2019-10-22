package daemon

import (
	"net"
	"os"
	"server/configure"
	"server/logger"

	"go.uber.org/zap"
)

// Start Server Running
func Start() {
	listener, err := net.Listen("tcp", configure.Single().TCP.Addr)
	if err != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, "connection error"); ce != nil {
			ce.Write(
				zap.Error(err),
			)
		}
		os.Exit(1)
	}
	if ce := logger.Logger.Check(zap.InfoLevel, "server work"); ce != nil {
		ce.Write(
			zap.String("addr", configure.Single().TCP.Addr),
		)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			if ce := logger.Logger.Check(zap.WarnLevel, "connection accpet error"); ce != nil {
				ce.Write(
					zap.Error(err),
				)
			}
			break
		}
		s := Session{
			Conn: conn,
		}
		go s.handConn()
	}
}
