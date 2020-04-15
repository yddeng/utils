package log

type LoggerI interface {
	Debugln(v ...interface{})
	Debugf(format string, v ...interface{})
	Infoln(v ...interface{})
	Infof(format string, v ...interface{})
	Warnln(v ...interface{})
	Warnf(format string, v ...interface{})
	Errorln(v ...interface{})
	Errorf(format string, v ...interface{})
	Fataln(v ...interface{})
	Fatalf(format string, v ...interface{})
}
