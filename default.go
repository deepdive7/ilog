package ilog

import "os"

var (
	gLogger Logger
)

func init() {
	gLogger = NewSimpleDefaultLogger(os.Stdout, 0, "Default->", true)
	gLogger.Deeper(1)
}

func SetDefaultLogger(logger Logger) {
	logger.Deeper(1)
	gLogger = logger
}

func SetLevel(l Level) {
	gLogger.SetLevel(l)
}

func SetPrefix(prefix string) {
	gLogger.SetPrefix(prefix)
}

func WithField(name string, value interface{}) error {
	return gLogger.WithField(name, value)
}

func WithFields(fields map[string]interface{}) error {
	return gLogger.WithFields(fields)
}

func AddReceiver(receiver *Receiver) {
	gLogger.AddReceiver(receiver)
}

func Debug(a ...interface{}) {
	gLogger.Debug(a...)
}

func Info(a ...interface{}) {
	gLogger.Info(a...)
}

func Warn(a ...interface{}) {
	gLogger.Warn(a...)
}

func Error(a ...interface{}) {
	gLogger.Error(a...)
}

func Fatal(a ...interface{}) {
	gLogger.Fatal(a...)
}

func Panic(a ...interface{}) {
	gLogger.Panic(a...)
}

func Debugf(format string, a ...interface{}) {
	gLogger.Debugf(format, a...)
}

func Infof(format string, a ...interface{}) {
	gLogger.Infof(format, a...)
}

func Warnf(format string, a ...interface{}) {
	gLogger.Warnf(format, a...)
}

func Errorf(format string, a ...interface{}) {
	gLogger.Errorf(format, a...)
}

func Fatalf(format string, a ...interface{}) {
	gLogger.Fatalf(format, a...)
}

func Panicf(format string, a ...interface{}) {
	gLogger.Panicf(format, a...)
}
