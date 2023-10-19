package main

import (
	"testing"

	assertions "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type Tests struct {
	suite.Suite
}

var (
	assert *assertions.Assertions
)

func (suite *Tests) SetupTest() {
	assert = assertions.New(suite.T())
}

func (suite *Tests) BeforeTest(suiteName, testName string) {
}

// func (suite *Tests) AfterTest(suiteName, testName string) {)

func TestSuite(t *testing.T) {
	cfg = &config{}
	parseConfig()
	StartMonitoringServer()
	suite.Run(t, new(Tests))
}
