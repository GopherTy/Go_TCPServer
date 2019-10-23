package daemon

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"server/logger"
	"server/protocal"

	"go.uber.org/zap"
)

// Session 用户的登录信息
type Session struct {
	Conn net.Conn
	File *os.File
}

// 处理连接
func (s *Session) handConn() {

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

		err = s.deal(msg)
		if err != nil {
			break
		}
	}

	logger.Logger.Debug("close conn",
		zap.String("remote", s.Conn.RemoteAddr().String()),
	)
	s.Conn.Close()
}

// deal 处理从字节数组中获取到的指令进行相应的服务器的操作。
func (s *Session) deal(m protocal.Message) (err error) {

	h := m.GetHeader()
	cmd := h.GetCmd()
	logger.Logger.Debug("deal ",
		zap.Uint16("cmd", cmd),
	)
	switch cmd {
	case protocal.HandlerAdd:
		// add
		return s.add(m)
	case protocal.HandlerDec:
		// dec
		return s.dec(m)
		// upload
	case protocal.HandlerUpload:
		return s.upload(m)
	case protocal.HandlerUploadRW:
		return s.uploadRw(m)
	case protocal.HandlerUploadRWOK:
		return s.uploadRwOK(m)

	default:
		return s.cmdNotExistError()
	}
}

// add 加法操作
func (s *Session) add(m protocal.Message) (err error) {
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

// dec 减法操作
func (s *Session) dec(m protocal.Message) (err error) {
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

// upload 服务器上传文件操作
func (s *Session) upload(m protocal.Message) (err error) {
	var content string
	data := m.GetBody()
	err = json.Unmarshal(data, &content)
	if err != nil {
		err = s.checkError(err, "daemon deal upload err ", "create file fail ")
		if err != nil {
			return
		}
	}
	// 判断文件是否存在
	s.uploadFileIsExist(content)
	return err
}

// uploadRw 进行文件创建和写入操作
func (s *Session) uploadRw(m protocal.Message) (err error) {
	b := make([]byte, ((protocal.BufLen-protocal.HandlerHeaderLength)-6)*3/4)
	data := m.GetBody()
	err = json.Unmarshal(data, &b)
	if err != nil {
		err = s.checkError(err, "deal upload create file error", "deal upload create file fail")
		if err != nil {
			return
		}
	}
	if err != nil {
		fmt.Println("open error ")
	}
	_, err = s.File.Write(b)
	if err != nil {
		fmt.Println("Read error", err)
	}
	return
}

// uploadRwOK
func (s *Session) uploadRwOK(m protocal.Message) (err error) {
	err = s.File.Close()
	if err != nil {
		return
	}
	var str string
	data := m.GetBody()
	err = json.Unmarshal(data, &str)
	if err != nil {
		err = s.checkError(err, "data send over ", " daemon data send error ")
		if err != nil {
			return
		}
	}
	msg, err := protocal.NewMessageFromJSON(protocal.HandlerUploadRWOK, str)
	if err != nil {
		err = s.checkError(err, "send over msg decode error", "send over msg error")
		if err != nil {
			return err
		}
	}
	_, err = s.Conn.Write(msg.Data)
	return
}

// checkError 用于处理指令出现错误时，服务器的错误处理
func (s *Session) checkError(e error, decodeStr, sendStr string) (err error) {
	if ce := logger.Logger.Check(zap.WarnLevel, decodeStr); ce != nil {
		ce.Write(
			zap.Error(e),
		)
	}
	msg, err := protocal.NewMessageFromJSON(protocal.HandlerFail, sendStr)
	if err != nil {
		return
	}
	_, err = s.Conn.Write(msg.Data)
	return
}

//uploadFileIsExist 判断文件是否存在
func (s *Session) uploadFileIsExist(content string) (err error) {
	file, err := os.OpenFile(content, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil && os.IsNotExist(err) {
		file, err = os.Create(content)
		if err != nil {
			if ce := logger.Logger.Check(zap.WarnLevel, "daemon creat file error"); ce != nil {
				ce.Write(zap.Error(err))
			}
		}
		file.Close()
		file, err = os.OpenFile(content, os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			if ce := logger.Logger.Check(zap.WarnLevel, "daemon creat file error"); ce != nil {
				ce.Write(zap.Error(err))
			}
		}
		s.File = file
		msg, err := protocal.NewMessageFromJSON(protocal.HandlerUploadRW, "ok")
		if err != nil {
			err = s.checkError(err, "daemon deal upload err ", "upload interrupt")
			if err != nil {
				return err
			}
		}
		_, err = s.Conn.Write(msg.Data)
		if err != nil {
			return err
		}
	} else {
		file.Close()
		msg, err := protocal.NewMessageFromJSON(protocal.HandlerUploadExist, "file is exist ")
		if err != nil {
			err = s.checkError(err, "daemon deal upload err ", "upload interrupt")
			if err != nil {
				return err
			}
		}
		_, err = s.Conn.Write(msg.Data)
		if err != nil {
			return err
		}
	}
	return
}

// cmdNotExistError 指令不存在时的处理
func (s *Session) cmdNotExistError() (err error) {
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
