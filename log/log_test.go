package log

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"testing"

	"github.com/labstack/gommon/log"

	"github.com/m0t0k1ch1/gommon/internal/testutil"
)

type ExitError int

func (err ExitError) Error() string {
	return fmt.Sprintf("exited with code %v", int(err))
}

func init() {
	exit = func(code int) {
		panic(ExitError(code))
	}
}

func newTestLogger(w io.Writer) *Logger {
	l := New("test")
	l.SetOutput(w)
	l.SetLevel(log.DEBUG)

	return l
}

func TestPrint(t *testing.T) {
	buf := new(bytes.Buffer)

	l := newTestLogger(buf)

	l.Print("print")
	testLog(t, buf.String(), "-", "print")

	buf.Reset()

	l.Printf("print%s", "f")
	testLog(t, buf.String(), "-", "printf")

	buf.Reset()

	l.Printj(log.JSON{"print": "j"})
	testLog(t, buf.String(), "-", `{"print":"j"}`)
}

func TestDebug(t *testing.T) {
	buf := new(bytes.Buffer)

	l := newTestLogger(buf)

	l.Debug("debug")
	testLog(t, buf.String(), "DEBUG", "debug")

	buf.Reset()

	l.Debugf("debug%s", "f")
	testLog(t, buf.String(), "DEBUG", "debugf")

	buf.Reset()

	l.Debugj(log.JSON{"debug": "j"})
	testLog(t, buf.String(), "DEBUG", `{"debug":"j"}`)
}

func TestInfo(t *testing.T) {
	buf := new(bytes.Buffer)

	l := newTestLogger(buf)

	l.Info("info")
	testLog(t, buf.String(), "INFO", "info")

	buf.Reset()

	l.Infof("info%s", "f")
	testLog(t, buf.String(), "INFO", "infof")

	buf.Reset()

	l.Infoj(log.JSON{"info": "j"})
	testLog(t, buf.String(), "INFO", `{"info":"j"}`)
}

func TestWarn(t *testing.T) {
	buf := new(bytes.Buffer)

	l := newTestLogger(buf)

	l.Warn("warn")
	testLog(t, buf.String(), "WARN", "warn")

	buf.Reset()

	l.Warnf("warn%s", "f")
	testLog(t, buf.String(), "WARN", "warnf")

	buf.Reset()

	l.Warnj(log.JSON{"warn": "j"})
	testLog(t, buf.String(), "WARN", `{"warn":"j"}`)
}

func TestError(t *testing.T) {
	buf := new(bytes.Buffer)

	l := newTestLogger(buf)

	l.Error("error")
	testLog(t, buf.String(), "ERROR", "error")

	buf.Reset()

	l.Errorf("error%s", "f")
	testLog(t, buf.String(), "ERROR", "errorf")

	buf.Reset()

	l.Errorj(log.JSON{"error": "j"})
	testLog(t, buf.String(), "ERROR", `{"error":"j"}`)
}

func TestFatal(t *testing.T) {
	buf := new(bytes.Buffer)

	l := newTestLogger(buf)

	testExit(t, func() { l.Fatal("fatal") }, 1)
	testLog(t, buf.String(), "FATAL", "fatal")

	buf.Reset()

	testExit(t, func() { l.Fatalf("fatal%s", "f") }, 1)
	testLog(t, buf.String(), "FATAL", "fatalf")

	buf.Reset()

	testExit(t, func() { l.Fatalj(log.JSON{"fatal": "j"}) }, 1)
	testLog(t, buf.String(), "FATAL", `{"fatal":"j"}`)
}

func TestPanic(t *testing.T) {
	buf := new(bytes.Buffer)

	l := newTestLogger(buf)

	testPanic(t, func() { l.Panic("panic") })
	testLog(t, buf.String(), "PANIC", "panic")

	buf.Reset()

	testPanic(t, func() { l.Panicf("panic%s", "f") })
	testLog(t, buf.String(), "PANIC", "panicf")

	buf.Reset()

	testPanic(t, func() { l.Panicj(log.JSON{"panic": "j"}) })
	testLog(t, buf.String(), "PANIC", `{"panic":"j"}`)
}

func TestLevel(t *testing.T) {
	buf := new(bytes.Buffer)

	l := newTestLogger(buf)

	l.SetLevel(log.INFO)
	testutil.Equal(t, l.Level(), log.INFO)

	l.Debug("debug")
	testutil.Equal(t, buf.String(), "")

	l.Info("info")
	testLog(t, buf.String(), "INFO", "info")
}

func TestHeader(t *testing.T) {
	buf := new(bytes.Buffer)

	l := newTestLogger(buf)

	l.SetLevel(log.INFO)
	l.SetHeader("${prefix} | ${level} | ")

	l.Info("info")
	testutil.Equal(t, buf.String(), "test | INFO | info\n")
}

func TestConcurrent(t *testing.T) {
	buf := new(bytes.Buffer)

	l := newTestLogger(buf)
	l.SetHeader("")

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			l.Debugf("debugf%d", i)
			wg.Done()
		}(i)
	}
	wg.Wait()

	body := buf.String()
	for i := 0; i < 100; i++ {
		testutil.Contains(t, body, fmt.Sprintf("debugf%d\n", i))
	}
}

func testLog(t *testing.T, body, levelName, message string) {
	t.Helper()

	testutil.Contains(t, body, fmt.Sprintf("prefix:%s", "test"))
	testutil.Contains(t, body, fmt.Sprintf("level:%s", levelName))
	testutil.Contains(t, body, fmt.Sprintf("message:%s", message))
}

func testExit(t *testing.T, f func(), code int) {
	t.Helper()

	defer func() {
		err := recover()
		switch v := err.(type) {
		case ExitError:
			if int(v) != code {
				t.Errorf("expected to exit with code %v but %v", code, int(v))
			}
		default:
			t.Error(err)
		}
	}()

	f()

	t.Errorf("expected to exit, but not")
}

func testPanic(t *testing.T, f func()) {
	t.Helper()

	defer func() {
		recover()
	}()

	f()

	t.Errorf("expected to panic, but not")
}
