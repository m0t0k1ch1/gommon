package log

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/labstack/gommon/log"
	"github.com/m0t0k1ch1/gommon/internal/testutils"
)

func newTestLogger(w io.Writer) *Logger {
	l := New("test")
	l.SetOutput(w)

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

	l.Printj(log.JSON{"key": "value"})
	testLog(t, buf.String(), "-", `{"key":"value"}`)
}

func testLog(t *testing.T, s, levelName, message string) {
	t.Helper()

	testutils.Contains(t, s, fmt.Sprintf("prefix:%s", "test"))
	testutils.Contains(t, s, fmt.Sprintf("level:%s", levelName))
	testutils.Contains(t, s, fmt.Sprintf("message:%s", message))
}
