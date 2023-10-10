package libpack_logging

import (
	"errors"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/suite"
)

type LoggingTestSuite struct {
	suite.Suite
}

var (
	testsLogger *LogConfig
)

type stdoutCapture struct {
	oldStdout *os.File
	readPipe  *os.File
}

func (sc *stdoutCapture) StartCapture() {
	sc.oldStdout = os.Stdout
	sc.readPipe, os.Stdout, _ = os.Pipe()
}

func (sc *stdoutCapture) StopCapture() (string, error) {
	if sc.oldStdout == nil || sc.readPipe == nil {
		return "", errors.New("StartCapture not called before StopCapture on Stdout")
	}
	os.Stdout.Close()
	os.Stdout = sc.oldStdout
	bytes, err := io.ReadAll(sc.readPipe)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

type stderrCapture struct {
	oldStderr *os.File
	readPipe  *os.File
}

func (sc *stderrCapture) StartCapture() {
	sc.oldStderr = os.Stderr
	sc.readPipe, os.Stderr, _ = os.Pipe()
}

func (sc *stderrCapture) StopCapture() (string, error) {
	if sc.oldStderr == nil || sc.readPipe == nil {
		return "", errors.New("StartCapture not called before StopCapture on Stderr")
	}
	os.Stderr.Close()
	os.Stderr = sc.oldStderr
	bytes, err := io.ReadAll(sc.readPipe)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (suite *LoggingTestSuite) SetupTest() {
}

func TestLoggingTestSuite(t *testing.T) {
	suite.Run(t, new(LoggingTestSuite))
}

func (suite *LoggingTestSuite) TestLogConfig_AllHandlers() {
	type args struct {
		message string
	}
	tests := []struct {
		name           string
		args           args
		wantLevel      string
		wantMessage    string
		envMinLogLevel string
		loggerType     string
		stdOutExpect   bool
		stdErrExpect   bool
	}{
		{
			name:       "Test log: Error",
			loggerType: "Error",
			args: args{
				message: "This is a error message",
			},
			wantLevel:    "error",
			wantMessage:  "This is a error message",
			stdErrExpect: true,
			stdOutExpect: false,
		},
		{
			name:       "Test log: Warning",
			loggerType: "Warning",
			args: args{
				message: "This is a warning message",
			},
			wantLevel:      "warn",
			wantMessage:    "This is a warning message",
			stdErrExpect:   false,
			stdOutExpect:   true,
			envMinLogLevel: "info",
		},
		{
			name:       "Test log: Warning | Min level: Debug",
			loggerType: "Warning",
			args: args{
				message: "This is a warning message",
			},
			wantLevel:      "warn",
			wantMessage:    "This is a warning message",
			stdErrExpect:   false,
			stdOutExpect:   true,
			envMinLogLevel: "debug",
		},
		{
			name:       "Test log: Info",
			loggerType: "Info",
			args: args{
				message: "This is a info message",
			},
			wantLevel:    "info",
			wantMessage:  "This is a info message",
			stdErrExpect: false,
			stdOutExpect: true,
		},
		{
			name:       "Test log: Info | Min level: Warn",
			loggerType: "Info",
			args: args{
				message: "This is a info message",
			},
			wantLevel:      "",
			wantMessage:    "",
			stdErrExpect:   false,
			stdOutExpect:   false,
			envMinLogLevel: "warn",
		},
		{
			name:       "Test log: Warning | Min level: Warn",
			loggerType: "Warning",
			args: args{
				message: "This is a warning message",
			},
			wantLevel:      "warn",
			wantMessage:    "This is a warning message",
			stdErrExpect:   false,
			stdOutExpect:   true,
			envMinLogLevel: "warn",
		},
		{
			name:       "Test log: Warning | Min level: Error",
			loggerType: "Warning",
			args: args{
				message: "This is an error message",
			},
			wantLevel:      "",
			wantMessage:    "",
			stdErrExpect:   false,
			stdOutExpect:   false,
			envMinLogLevel: "error",
		},
		{
			name:       "Test log: Debug | Min level: Debug",
			loggerType: "Debug",
			args: args{
				message: "This is a debug message",
			},
			wantLevel:      "debug",
			wantMessage:    "This is a debug message",
			stdErrExpect:   false,
			stdOutExpect:   true,
			envMinLogLevel: "debug",
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			if tt.envMinLogLevel != "" {
				os.Setenv("LOG_LEVEL", tt.envMinLogLevel)
				defer os.Unsetenv("LOG_LEVEL")
			}
			testsLogger = NewLogger()

			captureStdout := stdoutCapture{}
			captureStdout.StartCapture()
			captureStderr := stderrCapture{}
			captureStderr.StartCapture()

			reflect.ValueOf(testsLogger).MethodByName(tt.loggerType).Call([]reflect.Value{reflect.ValueOf(tt.args.message)})

			stdoutOut, err := captureStdout.StopCapture()
			if err != nil {
				suite.T().Fatal(err)
			}

			stderrOut, err := captureStderr.StopCapture()
			if err != nil {
				suite.T().Fatal(err)
			}

			if tt.stdErrExpect && !tt.stdOutExpect {
				gotLvl, gotMsg, err := getResponseValues(stderrOut, "short_message")
				suite.NoError(err, "Failed in [STDERR]: "+tt.name)
				suite.Equal(tt.wantLevel, gotLvl, "Failed in [STDERR]: "+tt.name)
				suite.Equal(tt.wantMessage, gotMsg, "Failed in [STDERR]: "+tt.name)
				suite.Equal("", stdoutOut, "Failed in [STDERR]: "+tt.name)
			}
			if tt.stdOutExpect && !tt.stdErrExpect {
				gotLvl, gotMsg, err := getResponseValues(stdoutOut, "short_message")
				suite.NoError(err, "Failed in [STDOUT]: "+tt.name)
				suite.Equal(tt.wantLevel, gotLvl, "Failed in [STDOUT]: "+tt.name)
				suite.Equal(tt.wantMessage, gotMsg, "Failed in [STDOUT]: "+tt.name)
				suite.Equal("", stderrOut, "Failed in [STDOUT]: "+tt.name)
			}
			if !tt.stdErrExpect && !tt.stdOutExpect {
				suite.Equal("", stderrOut, "Failed in [NEITHER]: "+tt.name)
				suite.Equal("", stdoutOut, "Failed in [NEITHER]: "+tt.name)
			}
			os.Unsetenv("LOG_LEVEL")
		})
	}
}

func (suite *LoggingTestSuite) TestFullMessage() {
	type args struct {
		extraFields map[string]interface{}
		message     string
	}
	extraFields := make(map[string]interface{})
	extraFields["_full_message"] = "full message"

	tests := []struct {
		args           args
		name           string
		wantLevel      string
		wantMessage    string
		envMinLogLevel string
		loggerType     string
		stdOutExpect   bool
		stdErrExpect   bool
	}{
		{
			name:       "Test log: Error",
			loggerType: "Error",
			args: args{
				message:     "This is a error message",
				extraFields: extraFields,
			},
			wantLevel:    "error",
			wantMessage:  extraFields["_full_message"].(string),
			stdErrExpect: true,
			stdOutExpect: false,
		},
		{
			name:       "Test log: Info",
			loggerType: "Info",
			args: args{
				message:     "This is a info message",
				extraFields: extraFields,
			},
			wantMessage:  extraFields["_full_message"].(string),
			stdErrExpect: false,
			stdOutExpect: true,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			if tt.envMinLogLevel != "" {
				os.Setenv("LOG_LEVEL", tt.envMinLogLevel)
				defer os.Unsetenv("LOG_LEVEL")
			}
			testsLogger = NewLogger()

			captureStdout := stdoutCapture{}
			captureStdout.StartCapture()
			captureStderr := stderrCapture{}
			captureStderr.StartCapture()

			reflect.ValueOf(testsLogger).MethodByName(tt.loggerType).Call([]reflect.Value{
				reflect.ValueOf(tt.args.message),
				reflect.ValueOf(tt.args.extraFields),
			})

			stdoutOut, err := captureStdout.StopCapture()
			if err != nil {
				suite.T().Fatal(err)
			}

			stderrOut, err := captureStderr.StopCapture()
			if err != nil {
				suite.T().Fatal(err)
			}

			if tt.stdErrExpect && !tt.stdOutExpect {
				_, gotMsg, err := getResponseValues(stderrOut, "_full_message")
				suite.NoError(err, "Failed in [STDERR]: "+tt.name)
				suite.Equal(tt.wantMessage, gotMsg, "Failed in [STDERR]: "+tt.name)
			}
			if tt.stdOutExpect && !tt.stdErrExpect {
				_, gotMsg, err := getResponseValues(stdoutOut, "_full_message")
				suite.NoError(err, "Failed in [STDOUT]: "+tt.name)
				suite.Equal(tt.wantMessage, gotMsg, "Failed in [STDOUT]: "+tt.name)
			}
			os.Unsetenv("LOG_LEVEL")
		})
	}
}

func Test_getResponseValues(t *testing.T) {
	type args struct {
		sourceJson string
	}
	tests := []struct {
		name       string
		args       args
		wantGotLvl string
		wantGotMsg string
		wantErr    bool
	}{
		{
			name: "Test with json",
			args: args{
				sourceJson: `{"level": "debug", "short_message": "hello world"`,
			},
			wantGotLvl: "debug",
			wantGotMsg: "hello world",
			wantErr:    false,
		},
		{
			name: "Test with json, wrong message field",
			args: args{
				sourceJson: `{"level": "debug", "message": "hello world"`,
			},
			wantGotLvl: "debug",
			wantGotMsg: "",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGotLvl, gotGotMsg, err := getResponseValues(tt.args.sourceJson, "short_message")
			if (err != nil) != tt.wantErr {
				t.Errorf("getResponseValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotGotLvl != tt.wantGotLvl {
				t.Errorf("getResponseValues() gotGotLvl = %v, want %v", gotGotLvl, tt.wantGotLvl)
			}
			if gotGotMsg != tt.wantGotMsg {
				t.Errorf("getResponseValues() gotGotMsg = %v, want %v", gotGotMsg, tt.wantGotMsg)
			}
		})
	}
}

func getResponseValues(sourceJson string, key string) (gotLvl, gotMsg string, err error) {
	gotLvl, err = jsonparser.GetString([]byte(sourceJson), "level")
	if err != nil {
		return
	}
	gotMsg, err = jsonparser.GetString([]byte(sourceJson), key)
	return
}
