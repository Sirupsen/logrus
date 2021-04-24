package logrus_test

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLogger_LogFn(t *testing.T) {
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.WarnLevel)

	notCalled := 0
	log.InfoFn(func() []interface{} {
		notCalled++
		return []interface{}{
			"Hello",
		}
	})
	assert.Equal(t, 0, notCalled)

	called := 0
	log.ErrorFn(func() []interface{} {
		called++
		return []interface{}{
			"Oopsi",
		}
	})
	assert.Equal(t, 1, called)
}
