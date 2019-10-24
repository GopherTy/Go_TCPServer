package function

import (
	"encoding/json"
	"server/cmd/daemon/resource"
	"server/logger"
	"server/protocal"

	"go.uber.org/zap"
)

// Dec 服务器减法操作
func Dec(m protocal.Message, s *resource.Session) (err error) {
	data := m.GetBody()
	var items []int
	err = json.Unmarshal(data, &items)
	if err != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, "daemon deal add error"); ce != nil {
			ce.Write(
				zap.Error(err),
			)
		}
		err = nil
		msg, err := protocal.NewMessageFromJSON(protocal.HandlerFail, "daemon deal dec error")
		_, err = s.Conn.Write(msg.Data)
		if err != nil {
			return err
		}
		return err
	}
	sub := 0
	for _, v := range items {
		sub -= v
	}
	msg, err := protocal.NewMessageFromJSON(protocal.HandlerSuccess, sub)
	if err != nil {
		return
	}
	_, err = s.Conn.Write(msg.Data)
	if err != nil {
		return
	}
	return
}
