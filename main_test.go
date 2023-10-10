package main

import (
	"fmt"
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
	fmt.Println("BeforeTest")
	cfg = &config{}
	parseConfig()
}

// func (suite *Tests) AfterTest(suiteName, testName string) {)

func TestSuite(t *testing.T) {
	suite.Run(t, new(Tests))
}
