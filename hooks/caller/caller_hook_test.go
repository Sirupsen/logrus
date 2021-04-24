package caller

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type logJson struct {
	Level string
	Msg   string
	Time  time.Time
	Src   string
}

func TestCallerHook_ExportedInfo(t *testing.T) {
	var buffer bytes.Buffer
	// Caller hook
	logrus.AddHook(NewHook(&CallerHookOptions{
		Field: "src",
		Flags: log.Lshortfile,
	}))
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(&buffer)

	testCallerHookExportedInfo(t, &buffer, "Testing")
	testCallerHookExportedInfo(t, &buffer, "Testing 2")
	testCallerHookExportedInfo(t, &buffer, "Testing 3")

	// Test in goroutines
	count := 50
	var wg sync.WaitGroup
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(t *testing.T, wg *sync.WaitGroup, i int) {
			defer wg.Done()
			logrus.Info(fmt.Sprintf("Testing for race conditions %d", i))
		}(t, &wg, i)
	}
	// Wait for goroutines
	wg.Wait()
}

func testCallerHookExportedInfo(t *testing.T, buffer *bytes.Buffer, msg string) {
	var j logJson
	logrus.Info(msg)
	scanner := bufio.NewScanner(buffer)
	scanner.Scan()
	line := scanner.Text()
	assert.NoError(t, scanner.Err(), "Scanner error")
	json.Unmarshal([]byte(line), &j)
	assert.Equal(t, "caller_test.go:54", j.Src, "%v", j)
	assert.Equal(t, msg, j.Msg, "%v", j)
	assert.Equal(t, "info", j.Level, "%v", j)
}

func TestCallerHook_ExportedEntryInfo(t *testing.T) {
	var buffer bytes.Buffer
	// Caller hook
	logrus.AddHook(NewHook(&CallerHookOptions{
		Field: "src",
		Flags: log.Lshortfile,
	}))
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(&buffer)

	testCallerHookExportedEntryInfo(t, &buffer, "Testing")
	testCallerHookExportedEntryInfo(t, &buffer, "Testing 2")
	testCallerHookExportedEntryInfo(t, &buffer, "Testing 3")

	// Test in goroutines
	count := 50
	var wg sync.WaitGroup
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(t *testing.T, wg *sync.WaitGroup, i int) {
			defer wg.Done()
			logrus.Info(fmt.Sprintf("Testing for race conditions %d", i))
		}(t, &wg, i)
	}
	// Wait for goroutines
	wg.Wait()
}

func testCallerHookExportedEntryInfo(t *testing.T, buffer *bytes.Buffer, msg string) {
	var j logJson
	logrus.WithField("testcaller", "testcaller value").Info(msg)
	scanner := bufio.NewScanner(buffer)
	scanner.Scan()
	line := scanner.Text()
	assert.NoError(t, scanner.Err(), "Scanner error")
	json.Unmarshal([]byte(line), &j)
	assert.Equal(t, "caller_test.go:95", j.Src, "%v", j)
	assert.Equal(t, msg, j.Msg, "%v", j)
	assert.Equal(t, "info", j.Level, "%v", j)
}

func TestCallerHook_NewLoggerInfo(t *testing.T) {
	var buffer bytes.Buffer
	logr := logrus.New()
	// Caller hook
	logr.Hooks.Add(NewHook(&CallerHookOptions{
		Field: "src",
		Flags: log.Lshortfile,
	}))
	logr.Formatter = &logrus.JSONFormatter{}
	logr.Out = &buffer

	testCallerHookInfo(t, logr, "Testing")
	testCallerHookInfo(t, logr, "Testing 2")
	testCallerHookInfo(t, logr, "Testing 3")

	// Test in goroutines
	count := 50
	var wg sync.WaitGroup
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(t *testing.T, wg *sync.WaitGroup, logr *logrus.Logger, i int) {
			defer wg.Done()
			logr.Info(fmt.Sprintf("Testing for race conditions %d", i))
		}(t, &wg, logr, i)
	}
	// Wait for goroutines
	wg.Wait()
}

func TestCallerHook_LoggerInfo(t *testing.T) {
	var buffer bytes.Buffer
	logr := &logrus.Logger{
		Out:       &buffer,
		Formatter: new(logrus.JSONFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	}
	// Caller hook
	logr.Hooks.Add(NewHook(&CallerHookOptions{
		Field: "src",
		Flags: log.Lshortfile,
	}))

	testCallerHookInfo(t, logr, "Testing")
	testCallerHookInfo(t, logr, "Testing 2")
	testCallerHookInfo(t, logr, "Testing 3")

	// Test in goroutines
	count := 50
	var wg sync.WaitGroup
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(t *testing.T, wg *sync.WaitGroup, logr *logrus.Logger, i int) {
			defer wg.Done()
			logr.Info(fmt.Sprintf("Testing for race conditions %d", i))
		}(t, &wg, logr, i)
	}
	// Wait for goroutines
	wg.Wait()
}

func testCallerHookInfo(t *testing.T, logr *logrus.Logger, msg string) {
	var j logJson
	logr.Info(msg)
	scanner := bufio.NewScanner(logr.Out.(*bytes.Buffer))
	scanner.Scan()
	line := scanner.Text()
	assert.NoError(t, scanner.Err(), "Scanner error")
	json.Unmarshal([]byte(line), &j)
	assert.Equal(t, "caller_test.go:169", j.Src, "%v", j)
	assert.Equal(t, msg, j.Msg, "%v", j)
	assert.Equal(t, "info", j.Level, "%v", j)
}
