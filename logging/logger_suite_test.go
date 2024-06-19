package libpack_logger

import (
	"testing"

	assertions "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type LoggerTestSuite struct {
	suite.Suite
}

var (
	assert *assertions.Assertions
)

func (suite *LoggerTestSuite) BeforeTest(suiteName, testName string) {
}

func (suite *LoggerTestSuite) SetupTest() {
	assert = assertions.New(suite.T())
}

// TearDownTest is run after each test to clean up
func (suite *LoggerTestSuite) TearDownTest() {
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(LoggerTestSuite))
}
