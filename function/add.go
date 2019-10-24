package function

import (
	"encoding/json"
	"server/cmd/daemon/resource"
	"server/logger"
	"server/protocal"

	"go.uber.org/zap"
)

// Add 加法操作
func Add(m protocal.Message, s *resource.Session) (err error) {
	b := m.GetBody()
	var items []int
	err = json.Unmarshal(b, &items)
	if err != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, "daemon deal add error"); ce != nil {
			ce.Write(
				zap.Error(err),
			)
		}
		err = nil
		msg, err := protocal.NewMessageFromJSON(protocal.HandlerFail, "add flag error")
		if err != nil {
			return err
		}
		_, err = s.Conn.Write(msg.Data)
		return err
	}

	// 进行计算，并且将计算结果返回给客户端。
	sum := 0
	for _, item := range items {
		sum += item
	}
	msg, err := protocal.NewMessageFromJSON(protocal.HandlerSuccess, sum)
	if err != nil {
		return
	}
	_, err = s.Conn.Write(msg.Data)
	if err != nil {
		return
	}
	return
}
