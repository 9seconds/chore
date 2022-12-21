package argparse_test

import (
	"testing"

	"github.com/9seconds/chore/internal/argparse"
	"github.com/9seconds/chore/internal/config"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/stretchr/testify/suite"
)

type ParseTestSuite struct {
	suite.Suite
	testlib.CtxTestSuite

	params map[string]config.Parameter
	flags  map[string]bool
}

func (suite *ParseTestSuite) SetupSuite() {
	intParam, _ := config.NewInteger(false, nil)
	strParam, _ := config.NewString(false, nil)
	reqParam, _ := config.NewString(true, nil)

	suite.params = map[string]config.Parameter{
		"int": intParam,
		"str": strParam,
		"req": reqParam,
	}
	suite.flags = map[string]bool{
		"cleanup": true,
		"welcome": false,
	}
}

func (suite *ParseTestSuite) SetupTest() {
	suite.CtxTestSuite.Setup(suite.T())
}

func (suite *ParseTestSuite) TestNothing() {
	args, err := argparse.Parse(suite.Context(), nil, nil, nil)
	suite.NoError(err)
	suite.Empty(args.Parameters)
	suite.Empty(args.Flags)
	suite.Empty(args.Positional)
}

func (suite *ParseTestSuite) TestAbsentRequiredParameter() {
	_, err := argparse.Parse(suite.Context(), nil, suite.flags, suite.params)
	suite.ErrorContains(err, "required but value is not provided")
}

func (suite *ParseTestSuite) TestOnlyRequiredParameters() {
	args, err := argparse.Parse(
		suite.Context(),
		[]string{"req=1", "+cleanup"},
		suite.flags,
		suite.params)
	suite.NoError(err)
	suite.Empty(args.Positional)

	suite.Len(args.Parameters, 1)
	suite.Equal("1", args.Parameters["req"])

	suite.Len(args.Flags, 1)
	suite.True(args.Flags["cleanup"])
}

func (suite *ParseTestSuite) TestMissRequiredFlag() {
	_, err := argparse.Parse(
		suite.Context(),
		[]string{"req=1"},
		suite.flags,
		suite.params)
	suite.ErrorContains(err, "flag")
	suite.ErrorContains(err, "is required but value is not provided")
}

func (suite *ParseTestSuite) TestMissRequiredParameter() {
	_, err := argparse.Parse(
		suite.Context(),
		[]string{"-cleanup"},
		suite.flags,
		suite.params)
	suite.ErrorContains(err, "parameter")
	suite.ErrorContains(err, "is required but value is not provided")
}

func (suite *ParseTestSuite) TestParseRepeatedParameters() {
	args, err := argparse.Parse(
		suite.Context(),
		[]string{"-cleanup", "req=1", "req=2"},
		suite.flags,
		suite.params)
	suite.NoError(err)
	suite.Equal("2", args.Parameters["req"])
}

func (suite *ParseTestSuite) TestParseRepeatedFlags() {
	args, err := argparse.Parse(
		suite.Context(),
		[]string{"req=3", "+cleanup", "+cleanup", "-cleanup"},
		suite.flags,
		suite.params)
	suite.NoError(err)
	suite.False(args.Flags["cleanup"])
}

func (suite *ParseTestSuite) TestParseUnknownParameter() {
	_, err := argparse.Parse(
		suite.Context(),
		[]string{"req=3", "+cleanup", "xx=3"},
		suite.flags,
		suite.params)
	suite.ErrorContains(err, "unknown parameter")
}

func (suite *ParseTestSuite) TestParseUnknownFlag() {
	_, err := argparse.Parse(
		suite.Context(),
		[]string{"req=3", "+cleanup", "-1"},
		suite.flags,
		suite.params)
	suite.ErrorContains(err, "unknown flag")
}

func (suite *ParseTestSuite) TestParseParameterAfterPositional() {
	_, err := argparse.Parse(
		suite.Context(),
		[]string{"+cleanup", "1", "req=3"},
		suite.flags,
		suite.params)
	suite.ErrorContains(err, "unexpected parameter")
}

func (suite *ParseTestSuite) TestParseFlagAfterPositional() {
	_, err := argparse.Parse(
		suite.Context(),
		[]string{"+cleanup", "1", "+cleanup"},
		suite.flags,
		suite.params)
	suite.ErrorContains(err, "unexpected flag")
}

func (suite *ParseTestSuite) TestBig() {
	args, err := argparse.Parse(
		suite.Context(),
		[]string{
			"+welcome",
			"-cleanup",
			"req=3",
			"int=1",
			"1",
			"2",
			"3",
			":req=4",
			":-cleanup",
		},
		suite.flags,
		suite.params)
	suite.NoError(err)

	suite.Equal([]string{"1", "2", "3", "req=4", "-cleanup"}, args.Positional)

	suite.Len(args.Flags, 2)
	suite.True(args.Flags["welcome"])
	suite.False(args.Flags["cleanup"])

	suite.Len(args.Parameters, 2)
	suite.Equal("3", args.Parameters["req"])
	suite.Equal("1", args.Parameters["int"])
}

func TestParse(t *testing.T) {
	suite.Run(t, &ParseTestSuite{})
}
