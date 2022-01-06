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
	fatalLevel log.Lvl = 6
	panicLevel log.Lvl = 7

	defaultLevel  = log.INFO
	defaultHeader = "time:${time}\tprefix:${prefix}\tlevel:${level}\tfile:${file}\tline:${line}"
)

var (
	defaultOutput = os.Stdout

	exit = os.Exit
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
			"FATAL",
			"PANIC",
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

func (l *Logger) Debug(a ...interface{}) {
	l.log(log.DEBUG, a...)
}

func (l *Logger) Debugf(format string, a ...interface{}) {
	l.logf(log.DEBUG, format, a...)
}

func (l *Logger) Debugj(j log.JSON) {
	l.logj(log.DEBUG, j)
}

func (l *Logger) Info(a ...interface{}) {
	l.log(log.INFO, a...)
}

func (l *Logger) Infof(format string, a ...interface{}) {
	l.logf(log.INFO, format, a...)
}

func (l *Logger) Infoj(j log.JSON) {
	l.logj(log.INFO, j)
}

func (l *Logger) Warn(a ...interface{}) {
	l.log(log.WARN, a...)
}

func (l *Logger) Warnf(format string, a ...interface{}) {
	l.logf(log.WARN, format, a...)
}

func (l *Logger) Warnj(j log.JSON) {
	l.logj(log.WARN, j)
}

func (l *Logger) Error(a ...interface{}) {
	l.log(log.ERROR, a...)
}

func (l *Logger) Errorf(format string, a ...interface{}) {
	l.logf(log.ERROR, format, a...)
}

func (l *Logger) Errorj(j log.JSON) {
	l.logj(log.ERROR, j)
}

func (l *Logger) Fatal(a ...interface{}) {
	l.log(fatalLevel, a...)
	exit(1)
}

func (l *Logger) Fatalf(format string, a ...interface{}) {
	l.logf(fatalLevel, format, a...)
	exit(1)
}

func (l *Logger) Fatalj(j log.JSON) {
	l.logj(fatalLevel, j)
	exit(1)
}

func (l *Logger) Panic(a ...interface{}) {
	l.log(panicLevel, a...)
	panic(fmt.Sprint(a...))
}

func (l *Logger) Panicf(format string, a ...interface{}) {
	l.logf(panicLevel, format, a...)
	panic(fmt.Sprintf(format, a...))
}

func (l *Logger) Panicj(j log.JSON) {
	l.logj(panicLevel, j)

	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}

	panic(string(b))
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
		panic(err)
	}

	buf.WriteByte('\t')
	buf.WriteString("message:" + message)
	buf.WriteByte('\n')

	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.output.Write(buf.Bytes())
}
