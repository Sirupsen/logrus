package logrus_sentry

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/getsentry/raven-go"
)

var (
	severityMap = map[logrus.Level]raven.Severity{
		logrus.DebugLevel: raven.DEBUG,
		logrus.InfoLevel:  raven.INFO,
		logrus.WarnLevel:  raven.WARNING,
		logrus.ErrorLevel: raven.ERROR,
		logrus.FatalLevel: raven.FATAL,
		logrus.PanicLevel: raven.FATAL,
	}
)

func getAndDel(d logrus.Fields, key string) (string, bool) {
	var (
		ok  bool
		v   interface{}
		val string
	)
	if v, ok = d[key]; !ok {
		return "", false
	}

	if val, ok = v.(string); !ok {
		return "", false
	}
	delete(d, key)
	return val, true
}

func getAndDelRequest(d logrus.Fields, key string) (*http.Request, bool) {
	var (
		ok  bool
		v   interface{}
		req *http.Request
	)
	if v, ok = d[key]; !ok {
		return nil, false
	}
	if req, ok = v.(*http.Request); !ok || req == nil {
		return nil, false
	}
	delete(d, key)
	return req, true
}

// SentryHook delivers logs to a sentry server.
type SentryHook struct {
	// Timeout sets the time to wait for a delivery error from the sentry server.
	// If this is set to zero the server will not wait for any response and will
	// consider the message correctly sent
	Timeout                 time.Duration
	StacktraceConfiguration stacktraceConfiguration

	client *raven.Client
	levels []logrus.Level
}

// stacktraceConfiguration allows for configuring stacktraces
type stacktraceConfiguration struct {
	// whether stacktraces should be enabled
	Enable bool
	// the level at which to start capturing stacktraces
	Level logrus.Level
	// how many stack frames to skip before stacktrace starts recording
	Skip int
	// the number of lines to include around a stack frame for context
	Context int
	// the prefixes that will be matched against the stack frame.
	// if the stack frame's package matches one of these prefixes
	// sentry will identify the stack frame as "in_app"
	InAppPrefixes []string
}

// NewSentryHook creates a hook to be added to an instance of logger
// and initializes the raven client.
// This method sets the timeout to 100 milliseconds.
func NewSentryHook(DSN string, levels []logrus.Level) (*SentryHook, error) {
	client, err := raven.NewClient(DSN, nil)
	if err != nil {
		return nil, err
	}
	return &SentryHook{
		Timeout: 100 * time.Millisecond,
		StacktraceConfiguration: stacktraceConfiguration{
			Enable:        false,
			Level:         logrus.ErrorLevel,
			Skip:          5,
			Context:       0,
			InAppPrefixes: nil,
		},
		client: client,
		levels: levels,
	}, nil
}

// Called when an event should be sent to sentry
// Special fields that sentry uses to give more information to the server
// are extracted from entry.Data (if they are found)
// These fields are: logger, server_name and http_request
func (hook *SentryHook) Fire(entry *logrus.Entry) error {
	packet := &raven.Packet{
		Message:   entry.Message,
		Timestamp: raven.Timestamp(entry.Time),
		Level:     severityMap[entry.Level],
		Platform:  "go",
	}

	d := entry.Data

	if logger, ok := getAndDel(d, "logger"); ok {
		packet.Logger = logger
	}
	if serverName, ok := getAndDel(d, "server_name"); ok {
		packet.ServerName = serverName
	}
	if req, ok := getAndDelRequest(d, "http_request"); ok {
		packet.Interfaces = append(packet.Interfaces, raven.NewHttp(req))
	}
	stConfig := &hook.StacktraceConfiguration
	if stConfig.Enable && entry.Level <= stConfig.Level {
		currentStacktrace := raven.NewStacktrace(stConfig.Skip, stConfig.Context, stConfig.InAppPrefixes)
		packet.Interfaces = append(packet.Interfaces, currentStacktrace)
	}
	packet.Extra = map[string]interface{}(d)

	_, errCh := hook.client.Capture(packet, nil)
	timeout := hook.Timeout
	if timeout != 0 {
		timeoutCh := time.After(timeout)
		select {
		case err := <-errCh:
			return err
		case <-timeoutCh:
			return fmt.Errorf("no response from sentry server in %s", timeout)
		}
	}
	return nil
}

// Levels returns the available logging levels.
func (hook *SentryHook) Levels() []logrus.Level {
	return hook.levels
}
