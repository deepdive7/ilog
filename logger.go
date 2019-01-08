package ilog

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
)

type Level int

// Available log levels
const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

type Logger interface {
	SetLevel(lvl Level)
	Deeper(delta int)
	SetPrefix(prefix string)
	WithField(name string, value interface{}) error
	WithFields(fields map[string]interface{}) error
	AddReceiver(receiver *Receiver)
	Debug(a ...interface{})
	Info(a ...interface{})
	Warn(a ...interface{})
	Error(a ...interface{})
	Fatal(a ...interface{})
	Debugf(format string, a ...interface{})
	Infof(format string, a ...interface{})
	Warnf(format string, a ...interface{})
	Errorf(format string, a ...interface{})
	Fatalf(format string, a ...interface{})
	Panic(a ...interface{})
	Panicf(format string, a ...interface{})
}

// Log options for all levels
var (
	optDebug = &levelOptions{"DEBUG", DEBUG, 34}
	optInfo  = &levelOptions{"INFO ", INFO, 32}
	optWarn  = &levelOptions{"WARN ", WARN, 33}
	optError = &levelOptions{"ERROR", ERROR, 31}
	optFatal = &levelOptions{"FATAL", FATAL, 31}

	DefaultCallerDepth = 3
)

// CallerInfoWithStructName returns the caller info, default depth is 3,
// caller->logger.Info->logger.logAll->CallerInfoWithStructName
func CallerInfoWithStructName(deepth int) string {
	stack := make([]byte, 1024, 1024)
	n := runtime.Stack(stack, false)
	info := ""
	if n <= 1024 {
		info = string(stack)
	} else {
		info = string(stack[:n])
	}
	callerStack := strings.Split(info, "\n")[2*deepth+1 : 2*deepth+3]
	funcInfo := callerStack[0]
	fileInfo := callerStack[1]
	file := fileInfo[strings.LastIndex(fileInfo, "/")+1:]
	return fmt.Sprintf("[%s <%s>]", funcInfo[:strings.LastIndex(funcInfo, "(")], file[:strings.Index(file, " ")])
}

func CallerInfo(deepth int) string {
	pc, file, line, ok := runtime.Caller(deepth)
	if !ok {
		return ""
	}
	ff := runtime.FuncForPC(pc)
	fname := ff.Name()
	i := strings.LastIndex(fname, ".")
	return fmt.Sprintf("[%s.%d <%s>]", file, line, fname[i+1:])
}

/* Format: [/home/edenzhong/github/ifly/peer.go:120->(*IPeer).Cleanup]
func CallerInfoWithStructName(level ...int) string {
	l := 3
	if len(level) > 0 && level[0] >= 0 {
		l = level[0]
	}
	_, file, line, ok := runtime.Caller(l)
	if !ok {
		return ""
	}
	//f := runtime.FuncForPC(pc)
	stack := make([]byte, 1024, 1024)
	n := runtime.Stack(stack, false)
	info := ""
	if n <= 1024 {
		info = string(stack)
	} else {
		info = string(stack[:n])
	}
	//funcName := f.Name()[strings.LastIndex(f.Name(), ".")+1:]
	funcName := ""
	callerStack := strings.Split(info, "\n")[7]
	left := strings.Index(callerStack, "(")
	right := strings.LastIndex(callerStack, "(")
	if left == right {
		funcName = callerStack[:right]
		if mi := strings.Index(funcName, "main."); mi != -1 {
			funcName = funcName[mi+5:]
		} else {
			li := strings.LastIndex(funcName, "/")
			funcName = funcName[li+1:]
			funcName = funcName[strings.Index(funcName, ".")+1:]
		}
	} else {
		funcName = callerStack[left:right]
	}
	return fmt.Sprintf("[%s:%d->%s]", file, line, W*/
// Options to store the key, level and color code of a log

type levelOptions struct {
	Key   string
	Level Level
	Color int
}

// Itol converts an integer to a logo.Level
func Itol(level int) Level {
	switch level {
	case 0:
		return DEBUG
	case 1:
		return INFO
	case 2:
		return WARN
	case 3:
		return ERROR
	case 4:
		return FATAL
	default:
		return DEBUG
	}
}

// DefaultLogger holds all Receivers
type DefaultLogger struct {
	Receivers      []*Receiver
	Active         bool
	NeedStructName bool
	CallerDeepth   int
	fields         map[string]interface{}
	fieldsStr      string
}

// NewDefaultLogger returns a new DefaultLogger filled with given Receivers
// Output format: prefix + [file.line <struct.func>]+{fields}: log contents
func NewDefaultLogger(recs ...*Receiver) *DefaultLogger {
	l := &DefaultLogger{
		Active:         true, // Every gLogger is active by default
		Receivers:      recs,
		CallerDeepth:   3,
		fields:         map[string]interface{}{},
		fieldsStr:      "",
	}
	return l
}

// NewSimpleDefaultLogger returns a gLogger with one simple Receiver
func NewSimpleDefaultLogger(w io.Writer, lvl Level, prefix string, color bool) *DefaultLogger {
	l := &DefaultLogger{}
	r := NewReceiver(w, prefix)
	r.Color = color
	r.Level = lvl
	l.Receivers = []*Receiver{r}
	l.Active = true
	l.fields = map[string]interface{}{}
	l.CallerDeepth = 3
	return l
}

func (l *DefaultLogger) Deeper(delta int) {
	l.CallerDeepth += delta
}

// SetLevel sets the log level of ALL receivers
func (l *DefaultLogger) SetLevel(lvl Level) {
	for _, r := range l.Receivers {
		r.Level = lvl
	}
}

// SetPrefix sets the prefix of ALL receivers
func (l *DefaultLogger) SetPrefix(prefix string) {
	for _, r := range l.Receivers {
		r.SetPrefix(prefix)
	}
}

// WithField add a field pair
func (l *DefaultLogger) WithField(name string, value interface{}) error {
	l.fields[name] = value
	data, err := json.Marshal(l.fields)
	if err != nil {
		return err
	}
	l.fieldsStr = string(data)
	return nil
}

// WithFields add a set of field pair
func (l *DefaultLogger) WithFields(fields map[string]interface{}) error {
	for k, v := range fields {
		l.fields[k] = v
	}

	data, err := json.Marshal(l.fields)
	if err != nil {
		return err
	}
	l.fieldsStr = string(data)
	return nil
}

// AddReceiver add a receiver to the gLogger's receivers list
func (l *DefaultLogger) AddReceiver(receiver *Receiver) {
	l.Receivers = append(l.Receivers, receiver)
}

// Write to all Receivers
func (l *DefaultLogger) logAll(opt *levelOptions, s string) {
	// Skip everything if gLogger is disabled
	if !l.Active {
		return
	}
	callerInfo := ""
	if opt.Level != 1 {
		if l.NeedStructName {
			callerInfo = CallerInfoWithStructName(l.CallerDeepth)
		} else {
			callerInfo = CallerInfo(l.CallerDeepth)
		}
	}
	ss := fmt.Sprintf("%s-%s: %s", callerInfo, l.fieldsStr, s)
	// Log to all receivers
	for _, r := range l.Receivers {
		r.log(opt, ss)
	}
}

// Debug logs arguments
func (l *DefaultLogger) Debug(a ...interface{}) {
	l.logAll(optDebug, fmt.Sprint(a...))
}

// Info logs arguments
func (l *DefaultLogger) Info(a ...interface{}) {
	l.logAll(optInfo, fmt.Sprint(a...))
}

// Warn logs arguments
func (l *DefaultLogger) Warn(a ...interface{}) {
	l.logAll(optWarn, fmt.Sprint(a...))
}

// Error logs arguments
func (l *DefaultLogger) Error(a ...interface{}) {
	l.logAll(optError, fmt.Sprint(a...))
}

// Fatal logs arguments
func (l *DefaultLogger) Fatal(a ...interface{}) {
	l.logAll(optFatal, fmt.Sprint(a...))
	os.Exit(1)
}

// Panic logs arguments
func (l *DefaultLogger) Panic(a ...interface{}) {
	s := fmt.Sprint(a...)
	l.logAll(optError, s)
	panic(s)
}

// Debugf logs formated arguments
func (l *DefaultLogger) Debugf(format string, a ...interface{}) {
	l.logAll(optDebug, fmt.Sprintf(format, a...))
}

// Infof logs formated arguments
func (l *DefaultLogger) Infof(format string, a ...interface{}) {
	l.logAll(optInfo, fmt.Sprintf(format, a...))
}

// Warnf logs formated arguments
func (l *DefaultLogger) Warnf(format string, a ...interface{}) {
	l.logAll(optWarn, fmt.Sprintf(format, a...))
}

// Errorf logs formated arguments
func (l *DefaultLogger) Errorf(format string, a ...interface{}) {
	l.logAll(optError, fmt.Sprintf(format, a...))
}

// Fatalf logs formated arguments
func (l *DefaultLogger) Fatalf(format string, a ...interface{}) {
	l.logAll(optFatal, fmt.Sprintf(format, a...))
	os.Exit(1)
}

// Panicf logs formated arguments
func (l *DefaultLogger) Panicf(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	l.logAll(optError, s)
	panic(s)
}

// Open is a short function to open a file with needed options
func Open(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
}

// Receiver holds all receiver options
type Receiver struct {
	// DefaultLogger object from the builtin log package
	Logger *log.Logger
	Level  Level
	Color  bool
	Active bool
	Format string
}

// SetPrefix sets the prefix of the gLogger.
// If a prefix is set and no trailing space is written, write one
func (r *Receiver) SetPrefix(prefix string) {
	if prefix != "" && !strings.HasSuffix(prefix, " ") {
		prefix += " "
	}
	r.Logger.SetPrefix(prefix)
}

// Logs to the gLogger
func (r *Receiver) log(opt *levelOptions, s string) {
	// Don't do anything if not wanted
	if !r.Active || opt.Level < r.Level {
		return
	}
	// Pre- and suffix
	prefix := ""
	suffix := "\n"
	// Add colors if wanted
	if r.Color {
		prefix += fmt.Sprintf("\x1b[0;%sm", strconv.Itoa(opt.Color))
		suffix = "\x1b[0m" + suffix
	}
	// Print to the gLogger
	r.Logger.Printf(prefix+r.Format+suffix, opt.Key, s)
}

// NewReceiver returns a new Receiver object with a given Writer
// and sets default values
func NewReceiver(w io.Writer, prefix string) *Receiver {
	logger := log.New(w, "", log.LstdFlags)
	r := &Receiver{
		Logger: logger,
	}
	// Default options
	r.Active = true
	r.Level = INFO
	r.Format = "[%s] ▶ %s"
	r.SetPrefix(prefix)
	return r
}
