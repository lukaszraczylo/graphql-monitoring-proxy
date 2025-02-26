package libpack_logger

import (
	"bytes"
	"testing"

	assertions "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// LoggerAdditionalTestSuite extends testing for functions with low coverage
type LoggerAdditionalTestSuite struct {
	suite.Suite
	logger *Logger
	output *bytes.Buffer
	assert *assertions.Assertions
}

func (suite *LoggerAdditionalTestSuite) SetupTest() {
	suite.output = &bytes.Buffer{}
	suite.logger = New().SetOutput(suite.output).SetShowCaller(false)
	suite.assert = assertions.New(suite.T())
}

func TestLoggerAdditionalTestSuite(t *testing.T) {
	suite.Run(t, new(LoggerAdditionalTestSuite))
}

// Test GetLogLevel function
func (suite *LoggerAdditionalTestSuite) TestGetLogLevel() {
	tests := []struct {
		name     string
		level    string
		expected int
	}{
		{"debug level", "debug", LEVEL_DEBUG},
		{"info level", "info", LEVEL_INFO},
		{"warn level", "warn", LEVEL_WARN},
		{"error level", "error", LEVEL_ERROR},
		{"fatal level", "fatal", LEVEL_FATAL},
		{"uppercase level", "DEBUG", LEVEL_DEBUG},
		{"mixed case level", "WaRn", LEVEL_WARN},
		{"invalid level", "invalid", defaultMinLevel},
		{"empty level", "", defaultMinLevel},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := GetLogLevel(tt.level)
			suite.assert.Equal(tt.expected, result)
		})
	}
}

// Test SetFieldName function
func (suite *LoggerAdditionalTestSuite) TestSetFieldName() {
	// Save original field names
	originalFieldNames := make(map[string]string)
	for k, v := range fieldNames {
		originalFieldNames[k] = v
	}

	// Restore original field names after test
	defer func() {
		for k, v := range originalFieldNames {
			fieldNames[k] = v
		}
	}()

	// Test with custom field names
	customTimestampField := "time"
	customLevelField := "severity"
	customMessageField := "text"

	suite.logger.SetFieldName("timestamp", customTimestampField)
	suite.logger.SetFieldName("level", customLevelField)
	suite.logger.SetFieldName("message", customMessageField)

	// Verify field names were changed
	suite.assert.Equal(customTimestampField, fieldNames["timestamp"])
	suite.assert.Equal(customLevelField, fieldNames["level"])
	suite.assert.Equal(customMessageField, fieldNames["message"])

	// Test logging with custom field names
	suite.output.Reset()
	suite.logger.Info(&LogMessage{Message: "test custom fields"})
	output := suite.output.String()

	// Check if custom field names are used in the output
	suite.assert.Contains(output, customTimestampField)
	suite.assert.Contains(output, customLevelField)
	suite.assert.Contains(output, customMessageField)
	suite.assert.NotContains(output, "timestamp")
	suite.assert.NotContains(output, "level")
	suite.assert.NotContains(output, "message")
}

// Test SetShowCaller and getCaller functions
func (suite *LoggerAdditionalTestSuite) TestSetShowCaller() {
	// Make sure caller info is disabled
	suite.logger.SetShowCaller(false)

	// Test with caller info disabled
	suite.output.Reset()
	suite.logger.Info(&LogMessage{Message: "test without cal__ler"})
	output := suite.output.String()
	suite.assert.NotContains(output, "caller")

	// Test with caller info enabled
	suite.output.Reset()
	suite.logger.SetShowCaller(true)
	suite.logger.Info(&LogMessage{Message: "test with caller"})
	output = suite.output.String()
	suite.assert.Contains(output, "caller")

	// Verify the caller info format (file:line)
	suite.assert.Regexp(`"caller":"[^:]+:\d+"`, output)
}

// Test Warning function
func (suite *LoggerAdditionalTestSuite) TestWarning() {
	suite.output.Reset()
	msg := &LogMessage{Message: "test warning"}
	suite.logger.Warning(msg)
	output := suite.output.String()
	suite.assert.Contains(output, "warn")
	suite.assert.Contains(output, "test warning")
}

// Test Error function
func (suite *LoggerAdditionalTestSuite) TestError() {
	suite.output.Reset()
	msg := &LogMessage{Message: "test error"}
	suite.logger.Error(msg)
	output := suite.output.String()
	suite.assert.Contains(output, "error")
	suite.assert.Contains(output, "test error")
}

// Test Fatal function
func (suite *LoggerAdditionalTestSuite) TestFatal() {
	suite.output.Reset()
	msg := &LogMessage{Message: "test fatal"}
	suite.logger.Fatal(msg)
	output := suite.output.String()
	suite.assert.Contains(output, "fatal")
	suite.assert.Contains(output, "test fatal")
}

// Test Critical function without exiting
func (suite *LoggerAdditionalTestSuite) TestCritical() {
	// Safely intercept os.Exit call with proper synchronization
	exitMutex.Lock()
	originalOsExit := osExit

	var exitCode int
	osExit = func(code int) {
		exitCode = code
		// Don't actually exit
	}
	exitMutex.Unlock()

	// Ensure we restore the original osExit function
	defer func() {
		exitMutex.Lock()
		osExit = originalOsExit
		exitMutex.Unlock()
	}()

	suite.output.Reset()
	msg := &LogMessage{Message: "test critical"}
	suite.logger.Critical(msg)
	output := suite.output.String()

	suite.assert.Contains(output, "fatal")
	suite.assert.Contains(output, "test critical")
	suite.assert.Equal(1, exitCode)
}
