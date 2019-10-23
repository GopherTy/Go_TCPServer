package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"server/configure"
	"server/logger"
	"server/protocal"

	"go.uber.org/zap"
)

//Run 运行客户端进程
func Run() {
	conn, err := net.Dial("tcp", configure.Single().TCP.Addr)
	if err != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, "client connection error:"); ce != nil {
			ce.Write(
				zap.Error(err),
			)
		}
		os.Exit(1)
	}
	// 发送数据
	content := "a.txt"
	msg, err := protocal.NewMessageFromJSON(protocal.HandlerUpload, content)
	checkError(conn, err)
	_, err = conn.Write(msg.Data)
	checkError(conn, err)
	for {
		msg, err = protocal.Read(conn)
		if err != nil {
			fmt.Println("client Read error:", err)
		}
		header := msg.GetHeader()
		if header.GetCmd() == protocal.HandlerUploadRW {
			path := "./cmd/client/a.txt"
			file, err := os.Open(path)
			if err != nil {
				fmt.Println("upload file fail :", err)
			}
			buf := make([]byte, ((protocal.BufLen-protocal.HandlerHeaderLength)-6)*3/4)
			err = sendData(conn, buf, file)
		} else if header.GetCmd() == protocal.HandlerUploadRWOK {
			var msgSucc string
			err := json.Unmarshal(msg.GetBody(), &msgSucc)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(msgSucc)
		} else if header.GetCmd() == protocal.HandlerUploadExist {
			var msgExst string
			err := json.Unmarshal(msg.GetBody(), &msgExst)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(msgExst)
		} else {
			var msgError string
			err := json.Unmarshal(msg.GetBody(), &msgError)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(msgError)
		}
	}

}

// 检查连接后读取错误。
func checkError(c net.Conn, err error) {
	if err != nil {
		fmt.Println("connection error: ", err)
		c.Close()
		os.Exit(1)
	}
}

// sendData 发送大文件数据
func sendData(conn net.Conn, buf []byte, file *os.File) (err error) {
	for {
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				msg, err := protocal.NewMessageFromJSON(protocal.HandlerUploadRWOK, "upload success")
				if err != nil {
					fmt.Println("client send data error : ", err)
					return err
				}
				_, err = conn.Write(msg.Data)
				if err != nil {
					fmt.Println("client write data error:", err)
					return err
				}
				break
			} else {
				fmt.Println("write error ", err)
				conn.Close()
				os.Exit(-1)
				break
			}
		}
		msg, err := protocal.NewMessageFromJSON(protocal.HandlerUploadRW, buf[:n])
		if err != nil {
			fmt.Println("client wirte error :", err)
			break
		}
		_, err = conn.Write(msg.Data)
		if err != nil {
			fmt.Println("client write", err)
			break
		}
	}
	return
}
