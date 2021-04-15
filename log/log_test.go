package log

import (
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger("log", "testLog", 100)
	//dlog.CloseStdOut()
	//logger.AsyncOut()
	logger.SetOutLevel(DEBUG, INFO)

	logger.Info("infoln message", 1)
	logger.Infof("%s : %d", "infof message", 2)
	time.Sleep(time.Second)
	logger.Debug("Debugln message", 1)
	logger.Debugf("%s : %d", "Debugf message", 2)
	time.Sleep(time.Second)
	logger.Error("Errorln message", 1)
	logger.Errorf("%s : %d", "Errorf message", 2)

	select {}

}
