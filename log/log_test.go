package log_test

import (
	dlog "github.com/yddeng/dutil/log"
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {
	logger := dlog.NewLogger("log", "testLog", 100)
	//dlog.CloseStdOut()
	//logger.AsyncOut()
	logger.SetOutLevel(dlog.DEBUG, dlog.INFO)

	logger.Infoln("infoln message", 1)
	logger.Infof("%s : %d", "infof message", 2)
	time.Sleep(time.Second)
	logger.Debugln("Debugln message", 1)
	logger.Debugf("%s : %d", "Debugf message", 2)
	time.Sleep(time.Second)
	logger.Errorln("Errorln message", 1)
	logger.Errorf("%s : %d", "Errorf message", 2)

	select {}

}
