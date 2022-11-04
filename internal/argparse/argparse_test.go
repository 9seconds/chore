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
}

func (suite *ParseTestSuite) SetupTest() {
	suite.CtxTestSuite.Setup(suite.T())
}

func (suite *ParseTestSuite) TestNothing() {
	intParam, _ := config.NewInteger(false, nil)
	strParam, _ := config.NewString(false, nil)
	params := map[string]config.Parameter{
		"int": intParam,
		"str": strParam,
	}

	args, err := argparse.Parse(suite.Context(), params, nil)
	suite.NoError(err)
	suite.Empty(args.Keywords)
	suite.Empty(args.Positional)
}

func (suite *ParseTestSuite) TestAbsentRequiredParameter() {
	_, err := argparse.Parse(suite.Context(), suite.params, nil)
	suite.ErrorContains(err, "absent value for parameter")
}

func (suite *ParseTestSuite) TestOnlyRequiredParameter() {
	args, err := argparse.Parse(suite.Context(), suite.params, []string{"req=1"})
	suite.NoError(err)
	suite.Len(args.Keywords, 1)
	suite.Equal("1", args.Keywords["req"])
	suite.Empty(args.Positional)
}

func (suite *ParseTestSuite) TestParseParameters() {
	args, err := argparse.Parse(
		suite.Context(),
		suite.params,
		[]string{"req=1", "int=1", "str=xx"})

	suite.NoError(err)
	suite.Len(args.Keywords, 3)
	suite.Equal("1", args.Keywords["req"])
	suite.Equal("1", args.Keywords["int"])
	suite.Equal("xx", args.Keywords["str"])
	suite.Empty(args.Positional)
}

func (suite *ParseTestSuite) TestInvalidValue() {
	_, err := argparse.Parse(suite.Context(), suite.params, []string{"req=1", "int=xx"})
	suite.ErrorContains(err, "incorrect value int for parameter")
}

func (suite *ParseTestSuite) TestUnknownParameter() {
	_, err := argparse.Parse(suite.Context(), suite.params, []string{"req=1", "xx=xx"})
	suite.ErrorContains(err, "unknown parameter")
}

func (suite *ParseTestSuite) TestParameterWithoutSeparator() {
	_, err := argparse.Parse(suite.Context(), suite.params, []string{"xx"})
	suite.ErrorContains(err, "cannot find = separator")
}

func (suite *ParseTestSuite) TestOnlyPositionals() {
	args, err := argparse.Parse(
		suite.Context(),
		suite.params,
		[]string{"req=1", "--", "1", "2", "3"})

	suite.NoError(err)
	suite.Len(args.Keywords, 1)
	suite.Equal([]string{"1", "2", "3"}, args.Positional)
}

func (suite *ParseTestSuite) TestNoPositionals() {
	args, err := argparse.Parse(suite.Context(), suite.params, []string{"req=1", "--"})
	suite.NoError(err)
	suite.Len(args.Keywords, 1)
	suite.Empty(args.Positional)
}

func (suite *ParseTestSuite) TestMergeArguments() {
	args, err := argparse.Parse(
		suite.Context(),
		suite.params,
		[]string{"req=1", "req=xx yy", "req='xx", "req=3"})
	suite.NoError(err)
	suite.Len(args.Keywords, 1)
	suite.Equal("1 'xx yy' ''\"'\"'xx' 3", args.Keywords["req"])
}

func TestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &ParseTestSuite{})
}
