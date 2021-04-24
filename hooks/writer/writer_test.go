package writer

import (
	"bytes"
	"io/ioutil"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestDifferentLevelsGoToDifferentWriters(t *testing.T) {
	var a, b bytes.Buffer

	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
		DisableColors:    true,
	})
	log.SetOutput(ioutil.Discard) // Send all logs to nowhere by default

	log.AddHook(&Hook{
		Writer: &a,
		LogLevels: []log.Level{
			log.WarnLevel,
		},
	})
	log.AddHook(&Hook{ // Send info and debug logs to stdout
		Writer: &b,
		LogLevels: []log.Level{
			log.InfoLevel,
		},
	})
	log.Warn("send to a")
	log.Info("send to b")

	assert.Equal(t, a.String(), "level=warning msg=\"send to a\"\n")
	assert.Equal(t, b.String(), "level=info msg=\"send to b\"\n")
}

func TestDifferentFormattersToDifferentWritter(t *testing.T) {
	var a, b bytes.Buffer

	log.SetOutput(ioutil.Discard) // Send all logs to nowhere by default

	log.AddHook(&Hook{
		Writer: &a,
		LogLevels: []log.Level{
			log.WarnLevel,
		},
		Formatter: &log.TextFormatter{
			DisableTimestamp: true,
			DisableColors:    true,
			FieldMap: log.FieldMap{
				log.FieldKeyLevel: "@level",
				log.FieldKeyMsg:   "@message",
			},
		},
	})
	log.AddHook(&Hook{ // Send info and debug logs to stdout
		Writer: &b,
		LogLevels: []log.Level{
			log.InfoLevel,
		},
		Formatter: &log.JSONFormatter{
			DisableTimestamp: true,
			FieldMap: log.FieldMap{
				log.FieldKeyMsg: "message",
			},
		},
	})
	log.Warn("send to a")
	log.Info("send to b")

	assert.Equal(t, a.String(), "@level=warning @message=\"send to a\"\n")
	assert.Equal(t, b.String(), "{\"level\":\"info\",\"message\":\"send to b\"}\n")
}
