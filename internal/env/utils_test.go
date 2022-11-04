package env_test

import (
	"testing"

	"github.com/9seconds/chore/internal/env"
	"github.com/alecthomas/assert/v2"
	"github.com/stretchr/testify/suite"
)

type MakeValueTestSuite struct {
	suite.Suite
}

func (suite *MakeValueTestSuite) TestValue() {
	testNames := [][]string{
		{"", "value", "=value"},
		{"x", "y", "x=y"},
		{"x", "y z", "x=y z"},
		{"x", "", "x="},
	}

	for _, values := range testNames {
		name := values[0]
		value := values[1]
		expected := values[2]

		suite.T().Run(name+"->"+value, func(t *testing.T) {
			assert.Equal(t, expected, env.MakeValue(name, value))
		})
	}
}

type EncodeBytesTestSuite struct {
	suite.Suite
}

func (suite *EncodeBytesTestSuite) TestValue() {
	suite.Equal(
		"AQIDBAU",
		env.EncodeBytes([]byte{1, 2, 3, 4, 5}))
}

func TestMakeValue(t *testing.T) {
	suite.Run(t, &MakeValueTestSuite{})
}

func TestEncodeBytes(t *testing.T) {
	suite.Run(t, &EncodeBytesTestSuite{})
}
