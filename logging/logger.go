// Package libpack_logger provides structured JSON logging with configurable
// log levels, caller information, and automatic sensitive data redaction.
// Supports debug, info, warning, and error log levels.
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
	"sync/atomic"
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
//
// timeFormat, minLogLevel and showCaller are read on every log() call, so
// they are stored as atomics rather than behind mu: this keeps the hot log
// path lock-free for the scalar/string config while still being safe to
// mutate concurrently from SetTimeFormat/SetMinLogLevel/SetShowCaller. mu
// continues to protect only the output writer, which requires exclusive
// access while writing bytes.
type Logger struct {
	output      io.Writer
	timeFormat  atomic.Pointer[string]
	minLogLevel atomic.Int32
	showCaller  atomic.Bool
	mu          sync.Mutex // Mutex to protect concurrent access to output
}

// LogMessage represents a log message with optional pairs.
type LogMessage struct {
	Pairs   map[string]any
	Message string
}

// bufferPool is used to reuse bytes.Buffer for efficiency.
var bufferPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

// fieldNames allows customization of output field names.
//
// It is package-global by design: a single process-wide field-name mapping
// is applied consistently across every Logger instance (this is the
// behavior existing callers rely on, e.g. main.go configures field names on
// one Logger and expects that naming to apply process-wide). fieldNamesMu
// guards it so SetFieldName (writer) and log (reader) can run concurrently
// without a data race.
var fieldNames = map[string]string{
	"timestamp": "timestamp",
	"level":     "level",
	"message":   "message",
}

// fieldNamesMu protects concurrent access to the package-global fieldNames map.
var fieldNamesMu sync.RWMutex

// osExit is a variable to allow mocking os.Exit in tests
var osExit = os.Exit

// exitMutex ensures thread-safe access to osExit
var exitMutex sync.RWMutex

// New creates a new Logger with default settings.
func New() *Logger {
	l := &Logger{
		output: os.Stdout,
	}
	format := defaultTimeFormat
	l.timeFormat.Store(&format)
	l.minLogLevel.Store(int32(defaultMinLevel))
	l.showCaller.Store(defaultShowCaller)
	return l
}

// SetOutput sets the output destination for the logger.
func (l *Logger) SetOutput(output io.Writer) *Logger {
	l.mu.Lock()
	l.output = output
	l.mu.Unlock()
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
	l.timeFormat.Store(&format)
	return l
}

// SetMinLogLevel sets the minimum log level for the logger.
func (l *Logger) SetMinLogLevel(level int) *Logger {
	l.minLogLevel.Store(int32(level))
	return l
}

// SetFieldName allows customizing the field names in log output.
func (l *Logger) SetFieldName(field, name string) *Logger {
	fieldNamesMu.Lock()
	fieldNames[field] = name
	fieldNamesMu.Unlock()
	return l
}

// SetShowCaller enables or disables including the caller information in log output.
func (l *Logger) SetShowCaller(show bool) *Logger {
	l.showCaller.Store(show)
	return l
}

// shouldLog determines if the message should be logged based on the logger's minimum log level.
func (l *Logger) shouldLog(level int) bool {
	return level >= int(l.minLogLevel.Load())
}

// IsLevelEnabled reports whether the given level would be emitted by this logger.
// Useful to gate expensive log-field construction (map/slice allocations) behind a
// cheap level check when the log call would otherwise be dropped.
func (l *Logger) IsLevelEnabled(level int) bool {
	return level >= int(l.minLogLevel.Load())
}

// log writes the log message with the given level.
func (l *Logger) log(level int, m *LogMessage) {
	if m.Pairs == nil {
		m.Pairs = make(map[string]any)
	}

	fieldNamesMu.RLock()
	timestampField := fieldNames["timestamp"]
	levelField := fieldNames["level"]
	messageField := fieldNames["message"]
	fieldNamesMu.RUnlock()

	m.Pairs[timestampField] = time.Now().Format(*l.timeFormat.Load())
	m.Pairs[levelField] = levelNames[level]
	m.Pairs[messageField] = m.Message

	if l.showCaller.Load() {
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
	// Lock the mutex before writing to the output to prevent race conditions
	l.mu.Lock()
	_, err = l.output.Write(buffer.Bytes())
	l.mu.Unlock()

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
