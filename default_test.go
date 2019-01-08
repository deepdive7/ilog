package ilog_test

import (
	"os"
	"testing"

	"github.com/deepdive7/ilog"
)
var l = ilog.NewSimpleDefaultLogger(os.Stdout, 0, "ILOG->", true)
func TestCallerInfo(t *testing.T) {
	l := ilog.NewSimpleDefaultLogger(os.Stdout, 1, "gLogger->", true)
	l.Error("Error")
}

// 100000	     9483 ns/op - 3083 ns/op = CallerInfo时间(6400ns)
func BenchmarkLoggerError(b *testing.B) {
	logFile, _ := os.Create("1.log")
	defer logFile.Close()
	//os.Stdout = logFile
	l := ilog.NewSimpleDefaultLogger(logFile, 1, "gLogger->", true)
	for i := 0; i < b.N; i++ {
		//fmt.Println("NewSimpleDefaultLogger(logFile, 1NewSimpleDefaultLogger(logFile, 1")
		l.Error("Error")
	}
}

type P struct {}

func (p *P) ThrowError() {
	l.Error("Throw error from *P")
}

func TestILOG(t *testing.T) {
	info := "INFO"
	debug := "Debug"
	warn := "WARN"
	err := "ERROR"
	//fatal := "FATAL"
	l.Info(info)
	l.Debug(debug)
	l.Warn(warn)
	l.Error(err)
	(&P{}).ThrowError()
	//l.Fatal(fatal)
}

/*
// 500000	      3070 ns/op
func BenchmarkColorError(b *testing.B) {
	f, _ := os.Create("1.log")
	os.Stdout = f
	s := color.New(color.FgBlack, color.BgLightCyan)
	for i := 0; i < b.N; i++ {
		s.Println("Logger not created correctlyLogger not created correctly")
	}
}
*/
