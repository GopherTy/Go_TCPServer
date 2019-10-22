package daemon

import (
	"server/configure"
	"server/logger"

	"go.uber.org/zap"
)

// Run run as daemon
func Run() {
	cnf := configure.Single()
	logger.Logger.Info("daemon running",
		zap.String("level", cnf.Logger.Level),
	)
	Start()
}
func f() string {
	return "f"
}
