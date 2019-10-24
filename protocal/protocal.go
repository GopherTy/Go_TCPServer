package protocal

import (
	"encoding/binary"
)

// Header 数据包的头部
type Header struct {
	Head []byte
}

// Message 数据包
type Message struct {
	Data []byte
}

// GetLen 获取头部中数据长度
func (h Header) GetLen() (l uint16) {
	return binary.LittleEndian.Uint16(h.Head)
}

// GetCmd 获取头部中的指令
func (h Header) GetCmd() (l uint16) {
	return binary.LittleEndian.Uint16(h.Head[2:])
}

// GetHeader 获取数据包中头部
func (m Message) GetHeader() (h Header) {
	h.Head = m.Data[:4]
	return h
}

// GetBody 获取数据包中的数据
func (m Message) GetBody() []byte {
	return m.Data[4:]
}
