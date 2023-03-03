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
	_, err := argparse.Parse([]string{"a\xc5z"}, argparse.DefaultListDelimiter)
	suite.ErrorContains(err, "is not valid UTF-8 string")
}

func (suite *ParseTestSuite) TestUnexpectedFlag() {
	for _, testName := range []string{"+c", "-c"} {
		testName := testName

		suite.T().Run(testName, func(t *testing.T) {
			_, err := argparse.Parse(
				[]string{"arg", testName},
				argparse.DefaultListDelimiter)
			assert.ErrorContains(t, err, "unexpected flag")
		})
	}
}

func (suite *ParseTestSuite) TestIncorrectFlag() {
	for _, testValue := range []string{"-", "+"} {
		testValue := testValue

		suite.T().Run(testValue, func(t *testing.T) {
			_, err := argparse.Parse(
				[]string{testValue},
				argparse.DefaultListDelimiter)
			assert.ErrorContains(t, err, "incorrect flag")
		})
	}
}

func (suite *ParseTestSuite) TestUnexpectedParameter() {
	_, err := argparse.Parse([]string{"arg", "c=1"}, argparse.DefaultListDelimiter)
	suite.ErrorContains(err, "unexpected parameter")
}

func (suite *ParseTestSuite) TestIncorrectParameter() {
	_, err := argparse.Parse([]string{"=1"}, argparse.DefaultListDelimiter)
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
		"p=:v:d:e:",
		"arg1",
		"arg2",
		":-j",
		":k=v",
	}, argparse.DefaultListDelimiter)

	suite.NoError(err)
	suite.Equal(map[string]string{
		"c": "3",
		"k": "2",
		"p": "v:d:e:",
	}, parsed.Parameters)
	suite.Equal(map[string]string{
		"x": argparse.FlagTrue,
		"k": argparse.FlagFalse,
		"j": argparse.FlagFalse,
	}, parsed.Flags)
	suite.Equal([]string{"arg1", "arg2", "-j", "k=v"}, parsed.Positional)
}

func (suite *ParseTestSuite) TestExplicitPositionals() {
	suite.T().Run("empty", func(t *testing.T) {
		parsed, err := argparse.Parse([]string{}, argparse.DefaultListDelimiter)
		assert.NoError(t, err)
		assert.False(t, parsed.IsPositionalTime())
	})

	suite.T().Run("non-empty", func(t *testing.T) {
		parsed, err := argparse.Parse(
			[]string{"--", "a", ":--", "b"},
			argparse.DefaultListDelimiter)
		assert.NoError(t, err)
		assert.True(t, parsed.IsPositionalTime())
		assert.Equal(t, []string{"a", "--", "b"}, parsed.Positional)
	})
}

func TestParse(t *testing.T) {
	suite.Run(t, &ParseTestSuite{})
}
