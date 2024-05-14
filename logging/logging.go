package libpack_logging

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/gookit/goutil/envutil"
	"github.com/rs/zerolog"
)

type LogConfig struct {
	logger zerolog.Logger
}

var (
	baseLogger zerolog.Logger

	eventPool = sync.Pool{
		New: func() interface{} {
			return new(zerolog.Event)
		},
	}

	fieldMapPool = sync.Pool{
		New: func() interface{} {
			return make(map[string]interface{})
		},
	}
)

func init() {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.MessageFieldName = "short_message"
	zerolog.TimestampFieldName = "timestamp"
	zerolog.LevelFieldName = "level"
	zerolog.LevelFatalValue = "critical"

	baseLogger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	switch logLevel := envutil.Getenv("LOG_LEVEL", "info"); logLevel {
	case "debug":
		baseLogger = baseLogger.Level(zerolog.DebugLevel)
	case "warn":
		baseLogger = baseLogger.Level(zerolog.WarnLevel)
	case "error":
		baseLogger = baseLogger.Level(zerolog.ErrorLevel)
	default:
		baseLogger = baseLogger.Level(zerolog.InfoLevel)
	}
}

func NewLogger() *LogConfig {
	return &LogConfig{logger: baseLogger}
}

func (lw *LogConfig) log(w io.Writer, level zerolog.Level, message string, fields map[string]interface{}) {
	logger := lw.logger.Output(w)
	event := logger.WithLevel(level).CallerSkipFrame(3)

	for k, val := range fields {
		switch v := val.(type) {
		case string:
			event = event.Str(k, v)
		case int:
			event = event.Int(k, v)
		case float64:
			event = event.Float64(k, v)
		default:
			event = event.Interface(k, val)
		}
	}

	event.Msg(message)
}

func (lw *LogConfig) logWithLevel(level zerolog.Level, message string, fields map[string]interface{}) {
	if lw.logger.GetLevel() <= level {
		w := os.Stdout
		if level >= zerolog.ErrorLevel {
			w = os.Stderr
		}
		lw.log(w, level, message, fields)
	}
}

func (lw *LogConfig) Debug(message string, fields map[string]interface{}) {
	lw.logWithLevel(zerolog.DebugLevel, message, fields)
}

func (lw *LogConfig) Info(message string, fields map[string]interface{}) {
	lw.logWithLevel(zerolog.InfoLevel, message, fields)
}

func (lw *LogConfig) Warning(message string, fields map[string]interface{}) {
	lw.logWithLevel(zerolog.WarnLevel, message, fields)
}

func (lw *LogConfig) Error(message string, fields map[string]interface{}) {
	lw.logWithLevel(zerolog.ErrorLevel, message, fields)
}

func (lw *LogConfig) Critical(message string, fields map[string]interface{}) {
	lw.logWithLevel(zerolog.FatalLevel, message, fields)
	os.Exit(1)
}

// Helper function to get a new fields map from the pool
func getFieldsMap() map[string]interface{} {
	return fieldMapPool.Get().(map[string]interface{})
}

// Helper function to put a used fields map back into the pool
func putFieldsMap(fields map[string]interface{}) {
	for k := range fields {
		delete(fields, k)
	}
	fieldMapPool.Put(fields)
}
