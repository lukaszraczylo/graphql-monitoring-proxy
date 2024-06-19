package libpack_logger

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/goccy/go-json"
)

func captureStderr(f func()) string {
	originalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w
	f()
	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	os.Stderr = originalStderr
	return buf.String()
}

func captureStdOut(f func()) string {
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	os.Stdout = originalStdout
	return buf.String()
}

func (suite *LoggerTestSuite) Test_LogMessageString() {
	msg := &LogMessage{
		Message: "test message",
	}

	assert.Equal("test message", msg.String())
}

func callLoggerMethod(logger *Logger, methodName string, message *LogMessage) {
	// Get the method by name using reflection
	method := reflect.ValueOf(logger).MethodByName(methodName)
	if method.IsValid() {
		// Call the method with the message as an argument
		method.Call([]reflect.Value{reflect.ValueOf(message)})
	} else {
		fmt.Printf("Method %s does not exist on Logger\n", methodName)
	}
}

func (suite *LoggerTestSuite) Test_LogsLevelsPrint() {
	output := &bytes.Buffer{}
	logger := New().SetOutput(output)

	tests := []struct {
		name            string
		method          string
		loggerMinLevel  int
		messageLogLevel int
		message         string
		pairs           map[string]any
		wantOutput      bool // Whether we expect output to be written
	}{
		{
			name:            "Log: Debug, Level: Debug - no pairs",
			method:          "Debug",
			loggerMinLevel:  LEVEL_DEBUG,
			messageLogLevel: LEVEL_DEBUG,
			message:         "debug message",
			wantOutput:      true,
		},
		{
			name:            "Log: Info, Level: Info - one pair",
			method:          "Info",
			loggerMinLevel:  LEVEL_INFO,
			messageLogLevel: LEVEL_INFO,
			message:         "info message",
			pairs: map[string]any{
				"key": "value",
			},
			wantOutput: true,
		},
		{
			name:            "Log: Info, Level: Warn - with pairs",
			method:          "Info",
			loggerMinLevel:  LEVEL_WARN,
			messageLogLevel: LEVEL_INFO,
			message:         "warn message",
			pairs: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
			wantOutput: false,
		},
		{
			name:            "Log: Warn, Level: Info - with 500 pairs",
			method:          "Warn",
			loggerMinLevel:  LEVEL_INFO,
			messageLogLevel: LEVEL_WARN,
			message:         "warn message with 500 pairs",
			pairs: func() map[string]any {
				pairs := make(map[string]any)
				for i := 0; i < 500; i++ {
					pairs[fmt.Sprintf("key%d", i)] = fmt.Sprintf("value%d", i)
				}
				return pairs
			}(),
			wantOutput: true,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			msg := &LogMessage{
				Message: tt.message,
				Pairs:   tt.pairs,
			}
			output.Reset()

			// Set logger's minimum log level
			logger.SetMinLogLevel(tt.loggerMinLevel)
			fmt.Println("Logger min log level:", LevelNames[logger.minLogLevel])

			// Call the logging method
			callLoggerMethod(logger, tt.method, msg)

			logOutput := output.String()
			fmt.Println("Output:", logOutput)

			if tt.wantOutput {
				var loggedMessage map[string]any
				err := json.Unmarshal([]byte(logOutput), &loggedMessage)
				if err != nil {
					t.Fatalf("Error unmarshalling log message: %v\nLog output: %s", err, logOutput)
				}

				if !containsLogMessage(logOutput, tt.message) {
					t.Errorf("Expected log message %q, but got %q", tt.message, logOutput)
				}
				assert.Equal(LevelNames[tt.messageLogLevel], loggedMessage["level"])
				if tt.pairs != nil {
					for k, v := range tt.pairs {
						assert.Equal(v, loggedMessage[k])
					}
				}
			} else {
				assert.Equal("", logOutput)
			}
		})
	}
}

func containsLogMessage(logOutput, expectedMessage string) bool {
	return bytes.Contains([]byte(logOutput), []byte(expectedMessage))
}

func (suite *LoggerTestSuite) Test_SetFormat() {
	logger := New().SetFormat(time.RFC3339Nano)

	assert.Equal(time.RFC3339Nano, logger.format)
}

func (suite *LoggerTestSuite) Test_SetMinLogLevel() {
	logger := New().SetMinLogLevel(LEVEL_DEBUG)

	assert.Equal(LEVEL_DEBUG, logger.minLogLevel)
}

func (suite *LoggerTestSuite) Test_ShouldLog() {
	logger := New().SetMinLogLevel(LEVEL_WARN)

	assert.True(logger.shouldLog(LEVEL_WARN))
	assert.True(logger.shouldLog(LEVEL_ERROR))
	assert.False(logger.shouldLog(LEVEL_INFO))
	assert.False(logger.shouldLog(LEVEL_DEBUG))
}
