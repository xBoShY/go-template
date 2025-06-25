package logging

import (
	"io"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

// Level refers to the log logging level
type Level uint32

// Create a general Base logger
var (
	baseLogger Logger
)

const (
	// Panic Level level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	Panic Level = iota
	// Fatal Level level. Logs and then calls `os.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	Fatal
	// Error Level level. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	Error
	// Warn Level level. Non-critical entries that deserve eyes.
	Warn
	// Info Level level. General operational entries about what's going on inside the
	// application.
	Info
	// Debug Level level. Usually only enabled when debugging. Very verbose logging.
	Debug
)

const stackPrefix = "[Stack]"

var once sync.Once

// Init needs to be called to ensure our logging has been initialized
func Init() {
	once.Do(func() {
		// By default, log to stderr (logrus's default), only warnings and above.
		baseLogger = NewLogger()
		baseLogger.SetLevel(Warn)
	})
}

func init() {
	Init()
}

// Fields maps logrus fields
type Fields = logrus.Fields

type JSON = map[string]any

// Logger is the interface for loggers.
type Logger interface {
	// Print logs a message at level Info.
	Print(...any)
	Println(...any)
	Printf(string, ...any)
	Printj(JSON)

	// Debug logs a message at level Debug.
	Debug(...any)
	Debugln(...any)
	Debugf(string, ...any)
	Debugj(JSON)

	// Info logs a message at level Info.
	Info(...any)
	Infoln(...any)
	Infof(string, ...any)
	Infoj(JSON)

	// Warn logs a message at level Warn.
	Warn(...any)
	Warnln(...any)
	Warnf(string, ...any)
	Warnj(JSON)

	// Error logs a message at level Error.
	Error(...any)
	Errorln(...any)
	Errorf(string, ...any)
	Errorj(JSON)

	// Fatal logs a message at level Fatal.
	Fatal(...any)
	Fatalln(...any)
	Fatalf(string, ...any)
	Fatalj(JSON)

	// Panic logs a message at level Panic.
	Panic(...any)
	Panicln(...any)
	Panicf(string, ...any)
	Panicj(JSON)

	// Add one key-value to log
	With(key string, value any) Logger

	// WithFields logs a message with specific fields
	WithFields(Fields) Logger

	// Set the logging version (Info by default)
	SetLevel(Level)

	// Get the logging version
	GetLevel() Level

	// Sets the output target
	SetOutput(io.Writer)

	// Sets the logger to JSON Format
	SetJSONFormatter()

	// Sets the logger to Text Format
	SetTextFormatter()

	IsLevelEnabled(level Level) bool

	// source adds file, line and function fields to the event
	source() *logrus.Entry

	// Adds a hook to the logger
	AddHook(hook logrus.Hook)
}

type logger struct {
	entry *logrus.Entry
}

func (l logger) With(key string, value any) Logger {
	return logger{
		l.entry.WithField(key, value),
	}
}

func (l logger) Print(args ...any) {
	l.Info(args...)
}

func (l logger) Println(args ...any) {
	l.Infoln(args...)
}

func (l logger) Printf(format string, args ...any) {
	l.Infof(format, args...)
}

func (l logger) Printj(j JSON) {
	l.Infoj(j)
}

func (l logger) Debug(args ...any) {
	l.source().Debug(args...)
}

func (l logger) Debugln(args ...any) {
	l.source().Debugln(args...)
}

func (l logger) Debugf(format string, args ...any) {
	l.source().Debugf(format, args...)
}

func (l logger) Debugj(j JSON) {
	l.source().WithFields(j).Debug()
}

func (l logger) Info(args ...any) {
	l.source().Info(args...)
}

func (l logger) Infoln(args ...any) {
	l.source().Infoln(args...)
}

func (l logger) Infof(format string, args ...any) {
	l.source().Infof(format, args...)
}

func (l logger) Infoj(j JSON) {
	l.source().WithFields(j).Info()
}

func (l logger) Warn(args ...any) {
	l.source().Warn(args...)
}

func (l logger) Warnln(args ...any) {
	l.source().Warnln(args...)
}

func (l logger) Warnf(format string, args ...any) {
	l.source().Warnf(format, args...)
}

func (l logger) Warnj(j JSON) {
	l.source().WithFields(j).Warn()
}

func (l logger) Error(args ...any) {
	l.source().Errorln(stackPrefix, string(debug.Stack()))
	l.source().Error(args...)
}

func (l logger) Errorln(args ...any) {
	l.source().Errorln(stackPrefix, string(debug.Stack()))
	l.source().Errorln(args...)
}

func (l logger) Errorf(format string, args ...any) {
	l.source().Errorln(stackPrefix, string(debug.Stack()))
	l.source().Errorf(format, args...)
}

func (l logger) Errorj(j JSON) {
	l.source().WithFields(j).Error()
}

func (l logger) Fatal(args ...any) {
	l.source().Errorln(stackPrefix, string(debug.Stack()))
	l.source().Fatal(args...)
}

func (l logger) Fatalln(args ...any) {
	l.source().Errorln(stackPrefix, string(debug.Stack()))
	l.source().Fatalln(args...)
}

func (l logger) Fatalf(format string, args ...any) {
	l.source().Errorln(stackPrefix, string(debug.Stack()))
	l.source().Fatalf(format, args...)
}

func (l logger) Fatalj(j JSON) {
	l.source().Errorln(stackPrefix, string(debug.Stack()))
	l.source().WithFields(j).Fatal()
}

func (l logger) Panic(args ...any) {
	defer func() {
		if r := recover(); r != nil {
			panic(r)
		}
	}()
	l.source().Errorln(stackPrefix, string(debug.Stack()))
	l.source().Panic(args...)
}

func (l logger) Panicln(args ...any) {
	defer func() {
		if r := recover(); r != nil {
			panic(r)
		}
	}()
	l.source().Errorln(stackPrefix, string(debug.Stack()))
	l.source().Panicln(args...)
}

func (l logger) Panicf(format string, args ...any) {
	defer func() {
		if r := recover(); r != nil {
			panic(r)
		}
	}()
	l.source().Errorln(stackPrefix, string(debug.Stack()))
	l.source().Panicf(format, args...)
}

func (l logger) Panicj(j JSON) {
	defer func() {
		if r := recover(); r != nil {
			panic(r)
		}
	}()
	l.source().Errorln(stackPrefix, string(debug.Stack()))
	l.source().WithFields(j).Panic()
}

func (l logger) WithFields(fields Fields) Logger {
	return logger{
		l.source().WithFields(fields),
	}
}

func (l logger) GetLevel() (lvl Level) {
	return Level(l.entry.Logger.Level)
}

func (l logger) SetLevel(lvl Level) {
	l.entry.Logger.Level = logrus.Level(lvl)
}

func (l logger) IsLevelEnabled(level Level) bool {
	return l.entry.Logger.Level >= logrus.Level(level)
}

func (l logger) SetOutput(w io.Writer) {
	l.setOutput(w)
}

func (l logger) setOutput(w io.Writer) {
	l.entry.Logger.Out = w
}

func (l logger) getOutput() io.Writer {
	return l.entry.Logger.Out
}

func (l logger) SetJSONFormatter() {
	l.entry.Logger.Formatter = &logrus.JSONFormatter{TimestampFormat: "2006-01-02T15:04:05.000000Z07:00"}
}

func (l logger) SetTextFormatter() {
	l.entry.Logger.Formatter = &logrus.TextFormatter{TimestampFormat: "2006-01-02T15:04:05.000000Z07:00"}
}

func (l logger) source() *logrus.Entry {
	event := l.entry

	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "<???>"
		line = 1
		event = event.WithFields(logrus.Fields{
			"file": file,
			"line": line,
		})
	} else {
		// Add file name and number
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
		event = event.WithFields(logrus.Fields{
			"file": file,
			"line": line,
		})

		// Add function name if possible
		if function := runtime.FuncForPC(pc); function != nil {
			event = event.WithField("function", function.Name())
		}
	}
	return event
}

func (l logger) AddHook(hook logrus.Hook) {
	l.entry.Logger.Hooks.Add(hook)
}

// Base returns the default Logger logging to
func Base() Logger {
	return baseLogger
}

// NewLogger returns a new Logger logging to out.
func NewLogger() Logger {
	l := logrus.New()
	return NewWrappedLogger(l)
}

// NewWrappedLogger returns a new Logger that wraps an external logrus logger.
func NewWrappedLogger(l *logrus.Logger) Logger {
	out := logger{
		logrus.NewEntry(l),
	}
	formatter := out.entry.Logger.Formatter
	tf, ok := formatter.(*logrus.TextFormatter)
	if ok {
		tf.TimestampFormat = "2006-01-02T15:04:05.000000 -0700"
	}
	return out
}

// RegisterExitHandler registers a function to be called on exit by logrus
// Exit handling happens when logrus.Exit is called, which is called by logrus.Fatal
func RegisterExitHandler(handler func()) {
	logrus.RegisterExitHandler(handler)
}
