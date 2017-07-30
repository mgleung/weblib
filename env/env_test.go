package env

import (
	. "gopkg.in/check.v1"
	"os"
	"testing"
	"time"
)

func Test(t *testing.T) {
	TestingT(t)
}

type EnvSuite struct{}

var _ = Suite(&EnvSuite{})

func (s *EnvSuite) TestToEnvKey(c *C) {
	testFlag := "this-is-test"
	expected := "THIS_IS_TEST"
	test := toEnvKey(testFlag)
	c.Assert(test, Equals, expected)
}

func (s *EnvSuite) TestEnvVariablesDuration(c *C) {
	if err := os.Setenv("TEST_DURATION", "12"); err != nil {
		c.Fatal(err)
	}
	var testDuration time.Duration
	FlagOrEnvDuration(&testDuration, "test-duration", 0, "")
	c.Assert(testDuration, Equals, time.Duration(12))

	if err := os.Setenv("TEST_DURATION_FAIL", "stringgggg"); err != nil {
		c.Fatal(err)
	}
	FlagOrEnvDuration(&testDuration, "test-duration-fail", 34, "")
	c.Assert(testDuration, Equals, time.Duration(34))
}
