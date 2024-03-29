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
	for _, testName := range []string{"+c", "_c"} {
		testName := testName

		suite.T().Run(testName, func(t *testing.T) {
			_, err := argparse.Parse(
				[]string{"arg", testName})
			assert.ErrorContains(t, err, "unexpected flag")
		})
	}
}

func (suite *ParseTestSuite) TestIncorrectFlag() {
	for _, testValue := range []string{"_", "+"} {
		testValue := testValue

		suite.T().Run(testValue, func(t *testing.T) {
			_, err := argparse.Parse(
				[]string{testValue})
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
		"_k",
		"_k",
		"_j",
		"p=d",
		"k=2",
		"c=3",
		"p=v",
		"arg1",
		"arg2",
		":_j",
		":k=v",
	})

	suite.NoError(err)
	suite.Equal(map[string][]string{
		"c": {"1", "3"},
		"k": {"2"},
		"p": {"d", "v"},
	}, parsed.Parameters)
	suite.Equal(map[string]bool{
		"x": true,
		"k": false,
		"j": false,
	}, parsed.Flags)
	suite.Equal([]string{"arg1", "arg2", "_j", "k=v"}, parsed.Positional)
}

func (suite *ParseTestSuite) TestExplicitPositionals() {
	suite.T().Run("empty", func(t *testing.T) {
		parsed, err := argparse.Parse([]string{})
		assert.NoError(t, err)
		assert.False(t, parsed.IsPositionalTime())
	})

	suite.T().Run("non-empty", func(t *testing.T) {
		parsed, err := argparse.Parse(
			[]string{"_x", "--", "a", ":--", "b", "-c"})
		assert.NoError(t, err)
		assert.True(t, parsed.IsPositionalTime())
		assert.Equal(t, []string{"a", ":--", "b", "-c"}, parsed.Positional)
	})
}

func TestParse(t *testing.T) {
	suite.Run(t, &ParseTestSuite{})
}
