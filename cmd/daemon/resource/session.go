package resource

import (
	"net"
	"os"
)

// Session 服务器端存储的连接对象
type Session struct {
	Conn net.Conn
	File *os.File
}
