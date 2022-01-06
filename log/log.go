package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/valyala/fasttemplate"
)

const (
	defaultLevel = log.INFO
)

var (
	defaultOutput = os.Stdout
	defaultHeader = "time:${time}\tprefix:${prefix}\tlevel:${level}\tfile:${file}\tline:${line}"
)

type Logger struct {
	output     io.Writer
	prefix     string
	level      log.Lvl
	levelNames []string
	header     *fasttemplate.Template
	bufferPool sync.Pool
	mutex      sync.Mutex
}

func New(prefix string) *Logger {
	l := &Logger{
		output: defaultOutput,
		prefix: prefix,
		level:  defaultLevel,
		levelNames: []string{
			"-",
			"DEBUG",
			"INFO",
			"WARN",
			"ERROR",
			"",
			"PANIC",
			"FATAL",
		},
		bufferPool: sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}

	l.setHeader(defaultHeader)

	return l
}

func (l *Logger) Output() io.Writer {
	return l.output
}

func (l *Logger) SetOutput(w io.Writer) {
	l.output = w
}

func (l *Logger) Prefix() string {
	return l.prefix
}

func (l *Logger) SetPrefix(prefix string) {
	l.prefix = prefix
}

func (l *Logger) Level() log.Lvl {
	return l.level
}

func (l *Logger) SetLevel(level log.Lvl) {
	l.level = level
}

func (l *Logger) SetHeader(header string) {
	l.setHeader(header)
}

func (l *Logger) setHeader(header string) {
	l.header = fasttemplate.New(header, "${", "}")
}

func (l *Logger) Print(a ...interface{}) {
	l.log(0, a...)
}

func (l *Logger) Printf(format string, a ...interface{}) {
	l.logf(0, format, a...)
}

func (l *Logger) Printj(j log.JSON) {
	l.logj(0, j)
}

func (l *Logger) log(level log.Lvl, a ...interface{}) {
	l.write(level, fmt.Sprint(a...))
}

func (l *Logger) logf(level log.Lvl, format string, a ...interface{}) {
	l.write(level, fmt.Sprintf(format, a...))
}

func (l *Logger) logj(level log.Lvl, j log.JSON) {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}

	l.write(level, string(b))
}

func (l *Logger) write(level log.Lvl, message string) {
	if level != 0 && level < l.level {
		return
	}

	buf := l.bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer l.bufferPool.Put(buf)

	_, file, line, _ := runtime.Caller(3)

	if _, err := l.header.ExecuteFunc(buf, func(w io.Writer, tag string) (int, error) {
		switch tag {
		case "time":
			return w.Write([]byte(strconv.FormatInt(time.Now().Unix(), 10)))
		case "prefix":
			return w.Write([]byte(l.prefix))
		case "level":
			return w.Write([]byte(l.levelNames[level]))
		case "file":
			return w.Write([]byte(file))
		case "line":
			return w.Write([]byte(strconv.Itoa(line)))
		default:
			return 0, nil
		}
	}); err != nil {
		return
	}

	buf.WriteByte('\t')
	buf.WriteString("message:" + message)
	buf.WriteByte('\n')

	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.output.Write(buf.Bytes())
}
