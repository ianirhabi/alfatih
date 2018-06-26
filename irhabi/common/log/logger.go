package log

import (
	"os"
	"path"
	"runtime"

	"github.com/Sirupsen/logrus"
)

var (
	// Default io.Writter.
	Writer = os.Stderr
	// If DebugMode JsonFormatter will be used.
	DebugMode = false
	// Default Logger
	Log = New()
)

// New creating default log with loggrus
// and custom formating.
func New(prefix ...string) *logrus.Logger {
	l := logrus.Logger{
		Out:       os.Stderr,
		Formatter: NewFormater(DebugMode, prefix...),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	}

	if DebugMode {
		l.Level = logrus.DebugLevel
	}

	return &l
}

// WithFile show file & line of caller.
// this useable for debuging and tracing errors.
func WithFile() *logrus.Entry {
	_, file, line, _ := runtime.Caller(2)
	return Log.WithFields(logrus.Fields{
		"file": path.Base(file),
		"line": line,
	})
}

// Debug print an message to the io.writer
// only when application is on debug mode.
func Debug(args ...interface{}) {
	l := WithFile()
	if len(args) == 1 {
		l.Debug(args[0])
	} else {
		l.Debugf(args[0].(string), args[1:]...)
	}
}

// Info print an message to the io.writer
// General operational entries about what's going on inside the
// application.
func Info(args ...interface{}) {
	if len(args) == 1 {
		Log.Info(args[0])
	} else {
		Log.Infof(args[0].(string), args[1:]...)
	}
}

// Warning print an message to the io.writer
// General operational entries about some warning.
func Warning(args ...interface{}) {
	l := WithFile()
	if len(args) == 1 {
		l.Warning(args[0])
	} else {
		l.Warningf(args[0].(string), args[1:]...)
	}
}

// Error print an error to the io.writer
// will show a stacktrace if is traced was has value.
func Error(err error) {
	if err != nil {
		WithFile().Error(err.Error())
	}
}
