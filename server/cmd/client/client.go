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
			break
		}
		header := msg.GetHeader()
		if header.GetCmd() == protocal.HandlerUpload {
			var result string
			err := json.Unmarshal(msg.GetBody(), &result)
			if err != nil {
				fmt.Println(err)
				return
			}
			if result == "ok" {
				path := "./cmd/client/a.txt"
				file, err := os.Open(path)
				if err != nil {
					fmt.Println("upload file fail :", err)
				}
				buf := make([]byte, protocal.BufLen-protocal.HandlerHeaderLength)
				// for len(buf) == protocal.BufLen {
				// 	n, _ := io.ReadFull(file, buf)
				// 	msg, err := protocal.NewMessageFromJSON(protocal.HandlerUploadRW, buf[:n])
				// 	if err != nil {
				// 		fmt.Println("client wirte error :", err)
				// 	}
				// 	_, err = conn.Write(msg.Data)
				// 	if err != nil {
				// 		fmt.Println("client write", err)
				// 	}
				// 	buf = buf[n:]
				// }
				n, _ := io.ReadFull(file, buf)
				fmt.Println(n)
				msg, err := protocal.NewMessageFromJSON(protocal.HandlerUploadRW, buf[:n])
				if err != nil {
					fmt.Println("client wirte error :", err)
				}
				_, err = conn.Write(msg.Data)
				if err != nil {
					fmt.Println("client write", err)
				}
			}
			var msgCreate string
			err = json.Unmarshal(msg.GetBody(), &msgCreate)
			if err != nil {
				fmt.Println("client read", err)
				return
			}
			fmt.Println(msgCreate)
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

// uploadFileToServer 上传文件到服务器
func uploadFileToServer(c net.Conn, err error) {

}
