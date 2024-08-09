package config

import (
	"os"
	"time"

	"github.com/xops-infra/noop/log"
)

func InitLog(debug bool) {
	logFile := "/opt/goosefs-cli2api/app.log"
	// mkdir&file
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		os.MkdirAll("/opt/goosefs-cli2api", 0755)
		os.Create(logFile)
	}

	if debug {
		log.Default().WithHumanTime(time.Local).WithLevel(log.DebugLevel).WithFilename(logFile).Init()
		log.Debugf("debug mode")
	} else {
		log.Default().WithHumanTime(time.Local).WithLevel(log.InfoLevel).WithFilename(logFile).Init()
	}

	log.Infof("init log success")
}
