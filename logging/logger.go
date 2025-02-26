package libpack_logger

import (
	"bytes"
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
	LEVEL_DEBUG = iota
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
	LEVEL_FATAL
)

var levelNames = []string{
	"debug",
	"info",
	"warn",
	"error",
	"fatal",
}

const (
	defaultTimeFormat = time.RFC3339
	defaultMinLevel   = LEVEL_INFO
	defaultShowCaller = false
)

// Logger represents the logging object with configurations.
type Logger struct {
	output      io.Writer
	timeFormat  string
	minLogLevel int
	showCaller  bool
}

// LogMessage represents a log message with optional pairs.
type LogMessage struct {
	Pairs   map[string]interface{}
	Message string
}

// bufferPool is used to reuse bytes.Buffer for efficiency.
var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// fieldNames allows customization of output field names.
var fieldNames = map[string]string{
	"timestamp": "timestamp",
	"level":     "level",
	"message":   "message",
}

// osExit is a variable to allow mocking os.Exit in tests
var osExit = os.Exit

// exitMutex ensures thread-safe access to osExit
var exitMutex sync.RWMutex

// New creates a new Logger with default settings.
func New() *Logger {
	return &Logger{
		timeFormat:  defaultTimeFormat,
		minLogLevel: defaultMinLevel,
		output:      os.Stdout,
		showCaller:  defaultShowCaller,
	}
}

// SetOutput sets the output destination for the logger.
func (l *Logger) SetOutput(output io.Writer) *Logger {
	l.output = output
	return l
}

// GetLogLevel returns the log level integer corresponding to the given level name.
func GetLogLevel(level string) int {
	level = strings.ToLower(level)
	for i, name := range levelNames {
		if name == level {
			return i
		}
	}
	return defaultMinLevel
}

// SetTimeFormat sets the time format for the logger's timestamp field.
func (l *Logger) SetTimeFormat(format string) *Logger {
	l.timeFormat = format
	return l
}

// SetMinLogLevel sets the minimum log level for the logger.
func (l *Logger) SetMinLogLevel(level int) *Logger {
	l.minLogLevel = level
	return l
}

// SetFieldName allows customizing the field names in log output.
func (l *Logger) SetFieldName(field, name string) *Logger {
	fieldNames[field] = name
	return l
}

// SetShowCaller enables or disables including the caller information in log output.
func (l *Logger) SetShowCaller(show bool) *Logger {
	l.showCaller = show
	return l
}

// shouldLog determines if the message should be logged based on the logger's minimum log level.
func (l *Logger) shouldLog(level int) bool {
	return level >= l.minLogLevel
}

// log writes the log message with the given level.
func (l *Logger) log(level int, m *LogMessage) {
	if m.Pairs == nil {
		m.Pairs = make(map[string]interface{})
	}

	m.Pairs[fieldNames["timestamp"]] = time.Now().Format(l.timeFormat)
	m.Pairs[fieldNames["level"]] = levelNames[level]
	m.Pairs[fieldNames["message"]] = m.Message

	if l.showCaller {
		m.Pairs["caller"] = getCaller()
	}

	buffer := bufferPool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer bufferPool.Put(buffer)

	encoder := json.NewEncoder(buffer)
	err := encoder.Encode(m.Pairs)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error marshalling log message:", err)
		return
	}

	_, err = l.output.Write(buffer.Bytes())
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error writing log message:", err)
	}
}

// Debug logs a debug-level message.
func (l *Logger) Debug(m *LogMessage) {
	if l.shouldLog(LEVEL_DEBUG) {
		l.log(LEVEL_DEBUG, m)
	}
}

// Info logs an info-level message.
func (l *Logger) Info(m *LogMessage) {
	if l.shouldLog(LEVEL_INFO) {
		l.log(LEVEL_INFO, m)
	}
}

// Warn logs a warning-level message.
func (l *Logger) Warn(m *LogMessage) {
	if l.shouldLog(LEVEL_WARN) {
		l.log(LEVEL_WARN, m)
	}
}

// Warning is an alias for Warn.
func (l *Logger) Warning(m *LogMessage) {
	l.Warn(m)
}

// Error logs an error-level message.
func (l *Logger) Error(m *LogMessage) {
	if l.shouldLog(LEVEL_ERROR) {
		l.log(LEVEL_ERROR, m)
	}
}

// Fatal logs a fatal-level message.
func (l *Logger) Fatal(m *LogMessage) {
	if l.shouldLog(LEVEL_FATAL) {
		l.log(LEVEL_FATAL, m)
	}
}

// Critical logs a critical-level message and exits the application.
func (l *Logger) Critical(m *LogMessage) {
	l.Fatal(m)
	exitMutex.RLock()
	defer exitMutex.RUnlock()
	osExit(1)
}

// getCaller retrieves the file and line number of the caller.
func getCaller() string {
	// Skip 3 stack frames: getCaller -> log -> [Debug|Info|...]
	const depth = 3
	_, file, line, ok := runtime.Caller(depth)
	if !ok {
		return "unknown:0"
	}
	file = filepath.Base(file)
	return fmt.Sprintf("%s:%d", file, line)
}
