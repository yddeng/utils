package log

import (
	"log"
	"os"
	"testing"
)

func TestLogger(t *testing.T) {

	Info("infoln message", 1)
	Infof("%s : %d", "infof message", 2)

	Debug("Debugln message", 1)
	Debugf("%s : %d", "Debugf message", 2)

	Error("Errorln message", 1)
	Errorf("%s : %d", "Errorf message", 2)

	SetPrefix("Prefix")
	Info("message info\n")
	Infof("message infof\n")

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

func TestNew(t *testing.T) {
	//logger := NewLogger(".", "", LdefFlags)
	//
	//logger.Info("logger")

	log.Print("log print")

	f, _ := os.Create("./test.log")
	log.SetOutput(f)

	log.Println("file log print")
}
