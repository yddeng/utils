package log

var logger *Logger

func InitLogger(l *Logger) {
	logger = l
}

func Debug(v ...interface{}) {
	if logger != nil {
		logger.output(DEBUG, "", v...)
	}
}

func Debugf(format string, v ...interface{}) {
	if logger != nil {
		logger.output(DEBUG, format, v...)
	}
}

func Info(v ...interface{}) {
	if logger != nil {
		logger.output(INFO, "", v...)
	}
}

func Infof(format string, v ...interface{}) {
	if logger != nil {
		logger.output(INFO, format, v...)
	}
}

func Warn(v ...interface{}) {
	if logger != nil {
		logger.output(WARN, "", v...)
	}
}

func Warnf(format string, v ...interface{}) {
	if logger != nil {
		logger.output(WARN, format, v...)
	}
}

func Error(v ...interface{}) {
	if logger != nil {
		logger.output(ERROR, "", v...)
	}
}

func Errorf(format string, v ...interface{}) {
	if logger != nil {
		logger.output(ERROR, format, v...)
	}
}
