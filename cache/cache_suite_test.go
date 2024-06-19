package libpack_cache

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	assertions "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type Tests struct {
	suite.Suite
}

var (
	assert             *assertions.Assertions
	redisMockServer, _ = miniredis.Run()
)

func (suite *Tests) BeforeTest(suiteName, testName string) {
}

func (suite *Tests) SetupTest() {
	cacheStats = &CacheStats{}
	assert = assertions.New(suite.T())
}

// TearDownTest is run after each test to clean up
func (suite *Tests) TearDownTest() {
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Tests))
}
