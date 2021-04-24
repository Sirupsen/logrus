package logrus

import (
	"fmt"
	"log"
	"strings"
)

// Fields type, used to pass to `WithFields`.
type Fields map[string]interface{}

// Level type
type Level uint32

// Convert the Level to a string. E.g. PanicLevel becomes "panic".
func (level Level) String() string {
	if b, err := level.MarshalText(); err == nil {
		return string(b)
	} else {
		return "unknown"
	}
}

// ParseLevel takes a string level and returns the Logrus log level constant.
func ParseLevel(lvl string) (Level, error) {
	switch strings.ToLower(lvl) {
	case "panic":
		return PanicLevel, nil
	case "fatal":
		return FatalLevel, nil
	case "error":
		return ErrorLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	case "trace":
		return TraceLevel, nil
	}

	var l Level
	return l, fmt.Errorf("not a valid logrus Level: %q", lvl)
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (level *Level) UnmarshalText(text []byte) error {
	l, err := ParseLevel(string(text))
	if err != nil {
		return err
	}

	*level = l

	return nil
}

func (level Level) MarshalText() ([]byte, error) {
	switch level {
	case TraceLevel:
		return []byte("trace"), nil
	case DebugLevel:
		return []byte("debug"), nil
	case InfoLevel:
		return []byte("info"), nil
	case WarnLevel:
		return []byte("warning"), nil
	case ErrorLevel:
		return []byte("error"), nil
	case FatalLevel:
		return []byte("fatal"), nil
	case PanicLevel:
		return []byte("panic"), nil
	}

	return nil, fmt.Errorf("not a valid logrus level %d", level)
}

// AllLevels constant exposing all logging levels
var AllLevels = []Level{
	PanicLevel,
	FatalLevel,
	ErrorLevel,
	WarnLevel,
	InfoLevel,
	DebugLevel,
	TraceLevel,
}

// These are the different logging levels. You can set the logging level to log
// on your instance of logger, obtained with `logrus.New()`.
const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

// Won't compile if StdLogger can't be realized by a log.Logger
var (
	_ StdLogger = &log.Logger{}
	_ StdLogger = &Entry{}
	_ StdLogger = &Logger{}

	_ FieldLogger = &Entry{}
	_ FieldLogger = &Logger{}
)

// StdLogger is what your logrus-enabled library should take, that way
// it'll accept a stdlib logger and a logrus logger. There's no standard
// interface, this is the closest we get, unfortunately.
type StdLogger interface {
	PrintLogger
	FatalLogger
	PanicLogger
}

type PanicLogger interface {
	Panic(...interface{})
	Panicf(string, ...interface{})
	Panicln(...interface{})
}

type PrintLogger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})
}
type FatalLogger interface {
	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fatalln(...interface{})
}

type WithFieldLogger interface {
	WithField(key string, value interface{}) *Entry
	WithFields(fields Fields) *Entry
	WithError(err error) *Entry
}

type DebugLogger interface {
	Debugf(format string, args ...interface{})
	Debug(args ...interface{})
	Debugln(args ...interface{})
}

type InfoLogger interface {
	Infof(format string, args ...interface{})
	Info(args ...interface{})
	Infoln(args ...interface{})
}

type WarnLogger interface {
	Warnf(format string, args ...interface{})
	Warn(args ...interface{})
	Warnln(args ...interface{})
}

type WarningLogger interface {
	Warningf(format string, args ...interface{})
	Warning(args ...interface{})
	Warningln(args ...interface{})
}

type ErrorLogger interface {
	Errorf(format string, args ...interface{})
	Error(args ...interface{})
	Errorln(args ...interface{})
}

// The FieldLogger interface generalizes the Entry and Logger types
type FieldLogger interface {
	StdLogger
	WithFieldLogger
	DebugLogger
	InfoLogger
	WarnLogger
	WarningLogger
	ErrorLogger

	// IsDebugEnabled() bool
	// IsInfoEnabled() bool
	// IsWarnEnabled() bool
	// IsErrorEnabled() bool
	// IsFatalEnabled() bool
	// IsPanicEnabled() bool
}

type TraceLogger interface {
	Tracef(format string, args ...interface{})
	Trace(args ...interface{})
	Traceln(args ...interface{})
}

// Ext1FieldLogger (the first extension to FieldLogger) is superfluous, it is
// here for consistancy. Do not use. Use Logger or Entry instead.
type Ext1FieldLogger interface {
	FieldLogger
	TraceLogger
}
