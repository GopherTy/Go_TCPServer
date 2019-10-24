package daemon

import (
	"server/cmd/daemon/resource"
	"server/logger"
	"server/protocal"

	"go.uber.org/zap"
)

// 处理连接
func handConn(s *resource.Session) {
	logger.Logger.Debug("new conn",
		zap.String("remote", s.Conn.RemoteAddr().String()),
	)
	for {
		msg, err := protocal.Read(s.Conn)
		if err != nil {
			if ce := logger.Logger.Check(zap.WarnLevel, "daemon read data error"); ce != nil {
				ce.Write(
					zap.Error(err),
				)
			}
			break
		}
		err = deal(s, msg)
		if err != nil {
			break
		}
	}
	logger.Logger.Debug("close conn",
		zap.String("remote", s.Conn.RemoteAddr().String()),
	)
	if s.File != nil {
		s.File.Close()
	}
	s.Conn.Close()
}

// deal 处理从字节数组中获取到的指令进行相应的服务器的操作。
func deal(s *resource.Session, m protocal.Message) (err error) {
	h := m.GetHeader()
	cmd := h.GetCmd()
	logger.Logger.Debug("deal ",
		zap.Uint16("cmd", cmd),
	)
	handler, ok := _Handlers[cmd]
	if ok {
		return handler(m, s)
	}
	err = cmdNotExistError(s)
	return
}

// cmdNotExistError 指令不存在时的处理
func cmdNotExistError(s *resource.Session) (err error) {
	if ce := logger.Logger.Check(zap.WarnLevel, "cmd error"); ce != nil {
		ce.Write(
			zap.Error(err),
		)
	}
	msg, err := protocal.NewMessageFromJSON(protocal.HandlerFail, "cmd not exist error")
	if err != nil {

		return
	}
	_, err = s.Conn.Write(msg.Data)
	return
}
