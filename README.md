# ilog
A simple logger.

### Simple Usage
```
func TestILOG(t *testing.T) {
	l := ilog.NewSimpleDefaultLogger(os.Stdout, 0, "ILOG->", true)
	info := "INFO"
	debug := "Debug"
	warn := "WARN"
	err := "ERROR"
	fatal := "FATAL"
	l.Info(info)
	l.Debug(debug)
	l.Warn(warn)
	l.Error(err)
	l.Fatal(fatal)
}
```

> Output screenshot

![ilog](./ilog.png)