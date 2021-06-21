package log

import (
	"testing"
)

func TestLogger(t *testing.T) {

	Info("infoln message", 1)
	Infof("%s : %d", "infof message", 2)

	Debug("Debugln message", 1)
	Debugf("%s : %d", "Debugf message", 2)

	Error("Errorln message", 1)
	Errorf("%s : %d", "Errorf message", 2)

	SetOutput("./", "testLog", 100)

	Debug("file debug")
	logger := Default()
	logger.Info("default info")
	Stack("test stack")

	CloseDebug()
	Debug("closed debug")

	CloseStdOut()
	Info("file info")

	Fatal("fatal")
}
