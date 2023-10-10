package libpack_logging

import (
	"io"
	"os"
	"time"

	"github.com/gookit/goutil/envutil"
	"github.com/rs/zerolog"
)

type LogConfig struct {
	logger zerolog.Logger
}

var baseLogger zerolog.Logger

func init() {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.MessageFieldName = "short_message"
	zerolog.TimestampFieldName = "timestamp"
	zerolog.LevelFieldName = "level"
	zerolog.LevelFatalValue = "critical"
	baseLogger = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

func NewLogger() *LogConfig {
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

	return &LogConfig{logger: baseLogger}
}

func (lw *LogConfig) log(w io.Writer, level zerolog.Level, message string, v map[string]interface{}) {
	e := lw.logger.With().Logger()
	e = e.Output(w)
	event := e.WithLevel(level).CallerSkipFrame(3)
	for k, val := range v {
		switch v := val.(type) {
		case string:
			event.Str(k, v)
		case int:
			event.Int(k, v)
		case float64:
			event.Float64(k, v)
		default:
			event.Interface(k, val)
		}
	}
	event.Msg(message)
}

func (lw *LogConfig) Debug(message string, v ...map[string]interface{}) {
	lw.log(os.Stdout, zerolog.DebugLevel, message, mergeMaps(v))
}

func (lw *LogConfig) Info(message string, v ...map[string]interface{}) {
	lw.log(os.Stdout, zerolog.InfoLevel, message, mergeMaps(v))
}

func (lw *LogConfig) Warning(message string, v ...map[string]interface{}) {
	lw.log(os.Stdout, zerolog.WarnLevel, message, mergeMaps(v))
}

func (lw *LogConfig) Error(message string, v ...map[string]interface{}) {
	lw.log(os.Stderr, zerolog.ErrorLevel, message, mergeMaps(v))
}

func (lw *LogConfig) Critical(message string, v ...map[string]interface{}) {
	lw.log(os.Stderr, zerolog.FatalLevel, message, mergeMaps(v))
	os.Exit(1)
}

func mergeMaps(maps []map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}
