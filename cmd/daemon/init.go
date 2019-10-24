package daemon

import (
	"server/configure"
	"server/function"
	"server/logger"
	"server/protocal"

	"go.uber.org/zap"
)

// Run run as daemon
func Run() {
	cnf := configure.Single()
	logger.Logger.Info("daemon running",
		zap.String("level", cnf.Logger.Level),
	)
	RegisterHandler(protocal.HandlerAdd, function.Add)
	RegisterHandler(protocal.HandlerDec, function.Dec)
	RegisterHandler(protocal.HandlerUpload, function.Upload)
	RegisterHandler(protocal.HandlerUploadRW, function.UploadRw)
	RegisterHandler(protocal.HandlerUploadRWOK, function.UploadRwOK)
	Start()
}
func f() string {
	return "f"
}
