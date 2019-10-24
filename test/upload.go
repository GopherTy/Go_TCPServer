package test

// import (
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"os"
// 	"server/logger"
// 	"server/protocal"

// 	"go.uber.org/zap"
// )

// // Upload 上传操作
// func Upload(m protocal.Message, s *Session) (err error) {
// 	var content string
// 	data := m.GetBody()
// 	err = json.Unmarshal(data, &content)
// 	if err != nil {
// 		err = s.CheckCmdError(err, "daemon deal upload err ", "create file fail ")
// 		if err != nil {
// 			return
// 		}
// 	}
// 	// 判断文件是否存在
// 	uploadFileIsExist(content, s)
// 	return err
// }

// // uploadFileIsExist 判断要上传的文件是否存在
// func uploadFileIsExist(content string, s *Session) (err error) {
// 	file, err := os.OpenFile(content, os.O_RDWR|os.O_APPEND, 0666)
// 	if err != nil && os.IsNotExist(err) {
// 		file, err = os.Create(content)
// 		if err != nil {
// 			if ce := logger.Logger.Check(zap.WarnLevel, "daemon creat file error"); ce != nil {
// 				ce.Write(zap.Error(err))
// 			}
// 		}
// 		file.Close()
// 		file, err = os.OpenFile(content, os.O_RDWR|os.O_APPEND, 0666)
// 		if err != nil {
// 			if ce := logger.Logger.Check(zap.WarnLevel, "daemon creat file error"); ce != nil {
// 				ce.Write(zap.Error(err))
// 			}
// 		}
// 		s.File = file
// 		msg, err := protocal.NewMessageFromJSON(protocal.HandlerUploadRW, "ok")
// 		if err != nil {
// 			err = s.CheckCmdError(err, "daemon deal upload err ", "upload interrupt")
// 			if err != nil {
// 				return err
// 			}
// 		}
// 		_, err = s.Conn.Write(msg.Data)
// 		if err != nil {
// 			return err
// 		}
// 	} else {
// 		file.Close()
// 		msg, err := protocal.NewMessageFromJSON(protocal.HandlerUploadExist, "file is exist ")
// 		if err != nil {
// 			err = s.CheckCmdError(err, "daemon deal upload err ", "upload interrupt")
// 			if err != nil {
// 				return err
// 			}
// 		}
// 		_, err = s.Conn.Write(msg.Data)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return
// }

// // uploadRw 正在进行上传操作
// func uploadRw(m protocal.Message, s *Session) (err error) {
// 	if s.File == nil {
// 		return errors.New("file nil")
// 	}
// 	b := make([]byte, ((protocal.BufLen-protocal.HandlerHeaderLength)-6)*3/4)
// 	data := m.GetBody()
// 	err = json.Unmarshal(data, &b)
// 	if err != nil {
// 		err = s.CheckCmdError(err, "deal upload create file error", "deal upload create file fail")
// 		if err != nil {
// 			return
// 		}
// 	}
// 	if err != nil {
// 		fmt.Println("open error ")
// 	}
// 	_, err = s.File.Write(b)
// 	if err != nil {
// 		fmt.Println("Read error", err)
// 	}
// 	return
// }

// // uploadRwOK
// func uploadRwOK(m protocal.Message, s *Session) (err error) {
// 	if s.File == nil {
// 		return errors.New("send error")
// 	}
// 	err = s.File.Close()
// 	if err != nil {
// 		return
// 	}
// 	s.File = nil
// 	var str string
// 	data := m.GetBody()
// 	err = json.Unmarshal(data, &str)
// 	if err != nil {
// 		err = s.CheckCmdError(err, "data send over ", " daemon data send error ")
// 		if err != nil {
// 			return
// 		}
// 	}
// 	msg, err := protocal.NewMessageFromJSON(protocal.HandlerUploadRWOK, str)
// 	if err != nil {
// 		err = s.CheckCmdError(err, "send over msg decode error", "send over msg error")
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	_, err = s.Conn.Write(msg.Data)
// 	return
// }
