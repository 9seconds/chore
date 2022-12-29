package argparse_test

import (
	"testing"

	"github.com/9seconds/chore/internal/argparse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ParseTestSuite struct {
	suite.Suite
}

func (suite *ParseTestSuite) TestInvalidArgument() {
	_, err := argparse.Parse([]string{"a\xc5z"})
	suite.ErrorContains(err, "is not valid UTF-8 string")
}

func (suite *ParseTestSuite) TestUnexpectedFlag() {
	for _, testName := range []string{"+c", "-c"} {
		testName := testName

		suite.T().Run(testName, func(t *testing.T) {
			_, err := argparse.Parse([]string{"arg", testName})
			assert.ErrorContains(t, err, "unexpected flag")
		})
	}
}

func (suite *ParseTestSuite) TestIncorrectFlag() {
	for _, testValue := range []string{"-", "+"} {
		testValue := testValue

		suite.T().Run(testValue, func(t *testing.T) {
			_, err := argparse.Parse([]string{testValue})
			assert.ErrorContains(t, err, "incorrect flag")
		})
	}
}

func (suite *ParseTestSuite) TestUnexpectedParameter() {
	_, err := argparse.Parse([]string{"arg", "c=1"})
	suite.ErrorContains(err, "unexpected parameter")
}

func (suite *ParseTestSuite) TestIncorrectParameter() {
	_, err := argparse.Parse([]string{"=1"})
	suite.ErrorContains(err, "incorrect parameter")
}

func (suite *ParseTestSuite) TestMixed() {
	parsed, err := argparse.Parse([]string{
		"c=1",
		"+x",
		"+k",
		"-k",
		"-k",
		"-j",
		"k=2",
		"c=3",
		"arg1",
		"arg2",
		":-j",
		":k=v",
	})

	suite.NoError(err)
	suite.Equal(map[string]string{
		"c": "3",
		"k": "2",
	}, parsed.Parameters)
	suite.Equal(map[string]argparse.FlagValue{
		"x": argparse.FlagTrue,
		"k": argparse.FlagFalse,
		"j": argparse.FlagFalse,
	}, parsed.Flags)
	suite.Equal([]string{"arg1", "arg2", "-j", "k=v"}, parsed.Positional)
}

func TestParse(t *testing.T) {
	suite.Run(t, &ParseTestSuite{})
}
