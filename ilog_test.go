package ilog_test

import (
	"github.com/deepdive7/ilog"
	"os"
	"testing"
)

//log_basic_test-> 2019/01/08 11:14:17 [INFO ] ▶ Info
//log_basic_test-> 2019/01/08 11:14:17 [DEBUG] ▶ [/home/edenzhong/github/ilog/ilog_test.go.13 <TestDefaultLogger>]Debug
//log_basic_test-> 2019/01/08 11:14:17 [WARN ] ▶ [/home/edenzhong/github/ilog/ilog_test.go.14 <TestDefaultLogger>]Warn
//log_basic_test-> 2019/01/08 11:14:17 [ERROR] ▶ [/home/edenzhong/github/ilog/ilog_test.go.15 <TestDefaultLogger>]Error
//log_basic_test-> 2019/01/08 11:14:17 [FATAL] ▶ [/home/edenzhong/github/ilog/ilog_test.go.16 <TestDefaultLogger>]Fatal
func TestDefaultLogger(t *testing.T) {
	ilog.SetPrefix("log_basic_test->")
	ilog.SetLevel(ilog.DEBUG)
	ilog.Info("Info")
	ilog.Debug("Debug")
	ilog.Warn("Warn")
	ilog.Error("Error")
	//ilog.Fatal("Fatal")
}

type LogCase struct{}

func (*LogCase) LogTest() {
	ilog.Info("Info")
	ilog.Debug("Debug")
	ilog.Warn("Warn")
	ilog.Error("Error")
	//ilog.Fatal("Fatal")
}

//log_struct_test-> 2019/01/08 11:13:43 [INFO ] ▶ Info
//log_struct_test-> 2019/01/08 11:13:43 [DEBUG] ▶ [github.com/deepdive7/ilog_test.(*LogCase).LogTest <ilog_test.go:22>]Debug
//log_struct_test-> 2019/01/08 11:13:43 [WARN ] ▶ [github.com/deepdive7/ilog_test.(*LogCase).LogTest <ilog_test.go:23>]Warn
//log_struct_test-> 2019/01/08 11:13:43 [ERROR] ▶ [github.com/deepdive7/ilog_test.(*LogCase).LogTest <ilog_test.go:24>]Error
//log_struct_test-> 2019/01/08 11:13:43 [FATAL] ▶ [github.com/deepdive7/ilog_test.(*LogCase).LogTest <ilog_test.go:25>]Fatal
func TestDefaultLoggerWithStructName(t *testing.T) {
	r := ilog.NewReceiver(os.Stdout, "")
	r.Level = ilog.DEBUG
	r.Color = true

	l := ilog.NewDefaultLogger(r)
	l.NeedStructName = true
	ilog.SetDefaultLogger(l)
	ilog.SetPrefix("log_struct_test->")
	ilog.SetLevel(ilog.DEBUG)
	(&LogCase{}).LogTest()
}

//log_fields_test-> 2019/01/08 11:25:20 [INFO ] ▶ -: Info without fields
//log_fields_test-> 2019/01/08 11:25:20 [INFO ] ▶ -{"id":100,"system":"glaucus-etl"}: Info
//log_fields_test-> 2019/01/08 11:25:20 [ERROR] ▶ [/home/edenzhong/github/ilog/ilog_test.go.63 <TestDefaultLogger_WithFields>]-{"id":100,"system":"glaucus-etl"}: Error
func TestDefaultLogger_WithFields(t *testing.T) {
	ilog.SetPrefix("log_fields_test->")
	ilog.SetLevel(ilog.DEBUG)
	ilog.Info("Info without fields")
	ilog.WithField("system", "glaucus-etl")
	ilog.WithFields(map[string]interface{}{
		"id": 100,
	})
	ilog.Info("Info")
	ilog.Error("Error")
}
