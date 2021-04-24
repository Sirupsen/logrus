package caller

import (
	"fmt"
	"log"
	"path"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// refer import "log"
const (
	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
)

// CallerHook adds caller information to log entries.
//
// Sample usage:
// logrus.AddHook(caller.NewHook(&caller.CallerHookOptions{
// 	Field: "src",
// 	Flags: log.Lshortfile,
// }))
// logrus.SetFormatter(&logrus.JSONFormatter{})
// logrus.Info("Test log")
// // time="2018-05-10T00:00:00-00:00" level=info msg="Test log" file="main.go:66"

type CallerHook struct {
	CallerHookOptions *CallerHookOptions
}

// NewHook creates a new caller hook with options. If options are nil or unspecified, options.Field defaults to "src"
// and options.Flags defaults to log.Llongfile
func NewHook(options *CallerHookOptions) *CallerHook {
	//old
	// Set default caller field to "src"
	if options.Field == "" {
		options.Field = "src"
	}
	// new
	if options.FileAlias == "" {
		options.FileAlias = "file"
	}
	if options.LineAlias == "" {
		options.LineAlias = "line"
	}
	// Set default caller flag to Std logger log.Llongfile
	if options.Flags == 0 {
		// old
		//options.Flags = log.Llongfile
		// new
		options.Flags = log.Lshortfile
	}
	return &CallerHook{options}
}

// CallerHookOptions stores caller hook options
type CallerHookOptions struct {
	// new
	FileAlias  string //default:file
	EnableFile bool
	LineAlias  string //default:line
	EnableLine bool
	// old
	// Field to display caller info in
	Field         string
	DisabledField bool
	// Stores the flags
	Flags int

	//fileNmar
	//Line Nmae
}

// HasFlag returns true if the report caller options contains the specified flag
func (options *CallerHookOptions) HasFlag(flag int) bool {
	return options.Flags&flag != 0
}

func (hook *CallerHook) Fire(entry *logrus.Entry) error {

	// new
	// get caller file and line here, it won't be available inside the goroutine
	// 1 for the function that called us.
	file, line := getCallerIgnoringLogMulti(1)
	if hook.CallerHookOptions.HasFlag(log.Lshortfile) && !hook.CallerHookOptions.HasFlag(log.Llongfile) {
		file = path.Base(file)
	}
	if !hook.CallerHookOptions.DisabledField {
		entry.Data[hook.CallerHookOptions.Field] = fmt.Sprintf("%s:%d", file, line)
	}
	if hook.CallerHookOptions.EnableFile {
		entry.Data[hook.CallerHookOptions.FileAlias] = file
	}
	if hook.CallerHookOptions.EnableLine {
		entry.Data[hook.CallerHookOptions.LineAlias] = line
	}

	// old
	//entry.Data[hook.CallerHookOptions.Field] = hook.callerInfo(entry.CallerFrames() + 1) // add 1 for this frame
	//entry.Data[hook.CallerHookOptions.Field] = hook.callerInfo(1) // add 1 for this frame

	return nil
}

func (hook *CallerHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}

//old
func (hook *CallerHook) callerInfo(skipFrames int) string {
	// Follows output of Std logger
	_, file, line, ok := runtime.Caller(skipFrames)
	if !ok {
		file = "???"
		line = 0
	} else {
		// check flags
		if hook.CallerHookOptions.HasFlag(log.Lshortfile) && !hook.CallerHookOptions.HasFlag(log.Llongfile) {
			file = path.Base(file)
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}

// getCaller returns the filename and the line info of a function
// further down in the call stack.  Passing 0 in as callDepth would
// return info on the function calling getCallerIgnoringLog, 1 the
// parent function, and so on.  Any suffixes passed to getCaller are
// path fragments like "/pkg/log/log.go", and functions in the call
// stack from that file are ignored.
func getCaller(callDepth int, suffixesToIgnore ...string) (file string, line int) {
	// bump by 1 to ignore the getCaller (this) stackframe
	callDepth++
outer:
	for {
		var ok bool
		_, file, line, ok = runtime.Caller(callDepth)
		if !ok {
			file = "???"
			line = 0
			break
		}

		for _, s := range suffixesToIgnore {
			if strings.HasSuffix(file, s) {
				callDepth++
				continue outer
			}
		}
		break
	}
	return
}

//new
func getCallerIgnoringLogMulti(callDepth int) (string, int) {
	// the +1 is to ignore this (getCallerIgnoringLogMulti) frame
	return getCaller(callDepth+1, "logrus/hooks.go", "logrus/entry.go", "logrus/logger.go", "logrus/exported.go", "asm_amd64.s")
}
