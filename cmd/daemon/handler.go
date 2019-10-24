package daemon

import (
	"fmt"
	"server/cmd/daemon/resource"
	"server/protocal"
)

// Handler 定义函数类型
type Handler func(m protocal.Message, s *resource.Session) (err error)

var _Handlers = make(map[uint16]Handler)

// RegisterHandler ...
func RegisterHandler(cmd uint16, handler Handler) {
	if handler == nil {
		panic("handler nil")
	}
	if _, ok := _Handlers[cmd]; ok {
		panic(fmt.Sprintf("cmd already exist %v", cmd))
	}
	_Handlers[cmd] = handler
	return
}
