package protocal

// 服务器与客户端之间的指令及状态码

// BufLen 缓冲区大小
const BufLen = 1024

// HandlerAdd 服务器端加法指令
const HandlerAdd = 1

// HandlerDec 服务器端减法指令
const HandlerDec = 2

// HandlerUpload 服务器上传文件
const HandlerUpload = 100

// HandlerUploadRW 读写文件
const HandlerUploadRW = 101

// HandlerSuccess 服务器端成功状态码
const HandlerSuccess = 0

// HandlerFail 服务器端失败状态码
const HandlerFail = 500

// HandlerHeaderLength 自定义协议头部长度
const HandlerHeaderLength = 4
