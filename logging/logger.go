package libpack_logger

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/goccy/go-json"
)

const (
	_ = iota
	LEVEL_DEBUG
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
	LEVEL_FATAL
)

var LevelNames = [...]string{
	"none",
	"debug",
	"info",
	"warn",
	"error",
	"fatal",
}

const (
	defaultFormat     = time.RFC3339
	defaultMinLevel   = LEVEL_INFO
	defaultShowCaller = false
)

var defaultOutput = os.Stdout

type Logger struct {
	output      io.Writer
	format      string
	minLogLevel int
	showCaller  bool
}

type LogMessage struct {
	output  io.Writer
	Pairs   map[string]any
	Message string
}

func (m *LogMessage) String() string {
	return m.Message
}

var fieldNames = map[string]string{
	"timestamp": "timestamp",
	"level":     "level",
	"message":   "message",
}

func New() *Logger {
	return &Logger{
		format:      defaultFormat,
		minLogLevel: defaultMinLevel,
		output:      defaultOutput,
		showCaller:  defaultShowCaller,
	}
}

func (l *Logger) SetOutput(output io.Writer) *Logger {
	l.output = output
	return l
}

var bufferPool = sync.Pool{
	New: func() any {
		return bytes.NewBuffer(nil)
	},
}

var defaultPairs = make(map[string]any)

func GetLogLevel(level string) int {
	for i, name := range LevelNames {
		if name == strings.ToLower(level) {
			return i
		}
	}
	return defaultMinLevel
}

func (l *Logger) log(level int, m *LogMessage) {
	if m.Pairs == nil {
		m.Pairs = defaultPairs
	}

	m.Pairs[fieldNames["timestamp"]] = time.Now().Format(l.format)
	m.Pairs[fieldNames["level"]] = LevelNames[level]
	m.Pairs[fieldNames["message"]] = m.Message

	if l.showCaller {
		m.Pairs["caller"] = getCaller()
	}

	buffer := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(buffer)
	buffer.Reset()

	encoder := json.NewEncoder(buffer)
	if err := encoder.Encode(m.Pairs); err != nil {
		fmt.Println("Error marshalling log message:", err)
		return
	}

	// if not running in test - use stderr and stdout, otherwise - use logger's output setting
	if flag.Lookup("test.v") == nil {
		if level >= LEVEL_ERROR {
			m.output = os.Stderr
		} else {
			m.output = os.Stdout
		}
	} else {
		m.output = l.output
	}

	m.output.Write(buffer.Bytes())
	m.output.Write([]byte("\n"))
}

func (l *Logger) Debug(m *LogMessage) {
	if l.shouldLog(LEVEL_DEBUG) {
		l.log(LEVEL_DEBUG, m)
	}
}

func (l *Logger) Info(m *LogMessage) {
	if l.shouldLog(LEVEL_INFO) {
		l.log(LEVEL_INFO, m)
	}
}

func (l *Logger) Warn(m *LogMessage) {
	if l.shouldLog(LEVEL_WARN) {
		l.log(LEVEL_WARN, m)
	}
}

func (l *Logger) Warning(m *LogMessage) {
	l.Warn(m)
}

func (l *Logger) Error(m *LogMessage) {
	if l.shouldLog(LEVEL_ERROR) {
		l.log(LEVEL_ERROR, m)
	}
}

func (l *Logger) Fatal(m *LogMessage) {
	if l.shouldLog(LEVEL_FATAL) {
		l.log(LEVEL_FATAL, m)
		os.Exit(1)
	}
}

func (l *Logger) Critical(m *LogMessage) {
	l.Fatal(m)
}

func (l *Logger) shouldLog(level int) bool {
	return level >= l.minLogLevel
}

func (l *Logger) SetFormat(format string) *Logger {
	l.format = format
	return l
}

func (l *Logger) SetMinLogLevel(level int) *Logger {
	l.minLogLevel = level
	return l
}

func (l *Logger) SetFieldName(field, name string) *Logger {
	fieldNames[field] = name
	return l
}

func (l *Logger) SetShowCaller(show bool) *Logger {
	l.showCaller = show
	return l
}

func getCaller() string {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		return "unknown:0"
	}
	file = filepath.Base(file)
	return fmt.Sprintf("%s:%d", file, line)
}
