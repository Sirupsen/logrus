// +build !windows,!nacl,!plan9

package syslog

import (
	"fmt"
	//"log/syslog"
	"github.com/GolangResources/syslog/syslog"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	sg "log/syslog"
)

const (
	NUM_SYSLOG_WORKERS = 8
	CHANNEL_DEPTH = 8192
)

// SyslogHook to send logs via syslog.
type SyslogHook struct {
	Writer		*syslog.Writer
	SyslogNetwork	string
	SyslogRaddr	string
	Tag		string
	syschan		chan *SEntry
}

type SEntry struct {
	line		string
	level		logrus.Level
}

// Creates a hook to be added to an instance of logger. This is called with
// `hook, err := NewSyslogHook("udp", "localhost:514", syslog.LOG_DEBUG, "")`
// `if err == nil { log.Hooks.Add(hook) }`
func NewSyslogHook(network, raddr string, priority sg.Priority, tag string) (*SyslogHook, error) {
	w, err := syslog.Dial("ctcp" , raddr, priority, nil)
	if err != nil {
		return nil, err
	}
	hook := &SyslogHook{
		w,
		network,
		raddr,
		tag,
		make(chan *SEntry, CHANNEL_DEPTH),
	}
	for i := 0; i <= NUM_SYSLOG_WORKERS; i++ {
		go hook.worker(i)
	}
	return hook, err
}

func (hook *SyslogHook) worker(i int) {
	for entry := range hook.syschan {
		time.Sleep(1)
		func() error {
			switch entry.level {
			case logrus.PanicLevel:
				return hook.Writer.Crit(&entry.line, hook.Tag)
			case logrus.FatalLevel:
				return hook.Writer.Crit(&entry.line, hook.Tag)
			case logrus.ErrorLevel:
				return hook.Writer.Err(&entry.line, hook.Tag)
			case logrus.WarnLevel:
				return hook.Writer.Warning(&entry.line, hook.Tag)
			case logrus.InfoLevel:
				return hook.Writer.Info(&entry.line, hook.Tag)
			case logrus.DebugLevel, logrus.TraceLevel:
				return hook.Writer.Debug(&entry.line, hook.Tag)
			default:
				return nil
			}
		}()
	}
}

func (hook *SyslogHook) Fire(entry *logrus.Entry) error {
	var err error
	line, err := entry.String()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}
	sentry := &SEntry{
		line:	line,
		level:	entry.Level,
	}
	select {
	case hook.syschan <- sentry:
	default:
		err = fmt.Errorf("logrus syslog chan full")
	}
	return err
}

func (hook *SyslogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
