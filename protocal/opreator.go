package protocal

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"net"
)

// 数据的操作：包括封包，拆包以及读取数据。

// NewMessageFromJSON 将数据通过json序列化为二进制数据在网络中传输
func NewMessageFromJSON(cmd uint16, v interface{}) (msg Message, err error) {
	body, err := json.Marshal(v)
	if err != nil {
		return
	}
	msg, err = NewMessage(cmd, body)
	return
}

// NewMessage 服务器和客户端之间的封包过程
func NewMessage(cmd uint16, body []byte) (msg Message, err error) {
	// 判断自定义消息大小是否超过缓冲区大小
	bodyLen := len(body)
	if bodyLen+HandlerHeaderLength > BufLen {
		err = errors.New("out of BufLenMax length")
		return
	}

	// 进行封包向客户端响应
	buf := make([]byte, bodyLen+HandlerHeaderLength)
	binary.LittleEndian.PutUint16(buf, uint16(bodyLen))
	binary.LittleEndian.PutUint16(buf[2:], cmd)

	copy(buf[4:], body)

	msg = Message{
		Data: buf,
	}
	return
}

// Read 按照自定义协议中的数据格式读取数据
func Read(conn net.Conn) (msg Message, err error) {
	buf := make([]byte, BufLen)
	err = readData(conn, buf[:4])
	if err != nil {
		return
	}
	header := &Header{
		Head: buf[:4],
	}
	l := header.GetLen()
	if l > BufLen-4 {
		return msg, errors.New("out of buffer lenght")
	}
	err = readData(conn, buf[4:l+4])
	if err != nil {
		return
	}
	msg.Data = buf[:l+4]
	return msg, nil
}

// readData 将缓冲区中的数据全部读取出来
func readData(conn net.Conn, data []byte) (err error) {

	_, err = io.ReadFull(conn, data)
	// for len(data) != 0 {
	// 	n, err := conn.Read(data)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	data = data[n:]
	// }
	if err != nil {
		return
	}
	return nil
}
