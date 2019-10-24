package function

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"server/cmd/daemon/resource"
	"server/logger"
	"server/protocal"

	"go.uber.org/zap"
)

// Upload 上传操作
func Upload(m protocal.Message, s *resource.Session) (err error) {
	var content string
	data := m.GetBody()
	err = json.Unmarshal(data, &content)
	if err != nil {
		err = checkCmdError(s, err, "send over msg decode error", "send over msg error")
		if err != nil {
			return
		}
	}
	// 判断文件是否存在
	uploadFileIsExist(content, s)
	return err
}

// uploadFileIsExist 判断要上传的文件是否存在
func uploadFileIsExist(content string, s *resource.Session) (err error) {
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
			err = checkCmdError(s, err, "send over msg decode error", "send over msg error")
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
			err = checkCmdError(s, err, "send over msg decode error", "send over msg error")
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

// UploadRw 正在进行上传操作
func UploadRw(m protocal.Message, s *resource.Session) (err error) {
	if s.File == nil {
		return errors.New("file nil")
	}
	b := make([]byte, ((protocal.BufLen-protocal.HandlerHeaderLength)-6)*3/4)
	data := m.GetBody()
	err = json.Unmarshal(data, &b)
	if err != nil {
		err = checkCmdError(s, err, "send over msg decode error", "send over msg error")
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

// UploadRwOK ...
func UploadRwOK(m protocal.Message, s *resource.Session) (err error) {
	if s.File == nil {
		return errors.New("send error")
	}
	err = s.File.Close()
	if err != nil {
		return
	}
	s.File = nil
	var str string
	data := m.GetBody()
	err = json.Unmarshal(data, &str)
	if err != nil {
		err = checkCmdError(s, err, "send over msg decode error", "send over msg error")
		if err != nil {
			return
		}
	}
	msg, err := protocal.NewMessageFromJSON(protocal.HandlerUploadRWOK, str)
	if err != nil {
		err = checkCmdError(s, err, "send over msg decode error", "send over msg error")
		if err != nil {
			return err
		}
	}
	_, err = s.Conn.Write(msg.Data)
	return
}

// CheckCmdError 用于处理指令出现错误时，服务器的错误处理
func checkCmdError(s *resource.Session, e error, decodeStr, sendStr string) (err error) {
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
