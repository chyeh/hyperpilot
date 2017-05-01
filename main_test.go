package main

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type UnitTestMainSuite struct{}

var _ = Suite(&UnitTestMainSuite{})

func (suite *UnitTestMainSuite) TestRunLoad(c *C) {
	testCases := []*struct {
		inputCPU                      float64
		inputMem                      float64
		expectedFiftiethPercentile    int
		expectedNinetyFifthPercentile int
		expectedNinetyNinthPercentile int
	}{
		{500.0, 1024.0, 400, 500, 550},
		{400.0, 1024.0, 600, 750, 825},
		{500.0, 512.0, 400, 500, 550},
		{500.0, 502.0, 700, 1000, 1600},
		{300.0, 512.0, 800, 1000, 1100},
		{300.0, 502.0, 1100, 1500, 2150},
	}

	c.Log(runLoad(406.875, 512.0))

	for i, testCase := range testCases {
		c.Logf("case[%d]: cpu[%f] mem[%f]", i+1, testCase.inputCPU, testCase.inputMem)

		actualP := runLoad(testCase.inputCPU, testCase.inputMem)
		c.Assert(actualP.FiftiethPercentile, Equals, testCase.expectedFiftiethPercentile)
		c.Assert(actualP.NinetyFifthPercentile, Equals, testCase.expectedNinetyFifthPercentile)
		c.Assert(actualP.NinetyNinthPercentile, Equals, testCase.expectedNinetyNinthPercentile)
	}
}

type IntegrationTestMainSuite struct{}

var _ = Suite(&IntegrationTestMainSuite{})

func (suite *IntegrationTestMainSuite) TestMyAlgo(c *C) {
	myAlgo(Services)
}
