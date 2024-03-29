package argparse_test

import (
	"testing"

	"github.com/9seconds/chore/internal/argparse"
	"github.com/9seconds/chore/internal/script/config"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ParsedArgsTestSuite struct {
	suite.Suite
	testlib.CtxTestSuite

	params map[string]config.Parameter
	flags  map[string]config.Flag
}

func (suite *ParsedArgsTestSuite) SetupTest() {
	suite.CtxTestSuite.Setup(suite.T())

	param1, err := config.NewInteger("int1", false, nil)
	require.NoError(suite.T(), err)

	param2, err := config.NewJSON("json1", true, nil)
	require.NoError(suite.T(), err)

	suite.params = map[string]config.Parameter{
		"int1":  param1,
		"json1": param2,
	}
	suite.flags = map[string]config.Flag{
		"flag1": config.NewFlag("flag1", true),
		"flag2": config.NewFlag("flag2", false),
	}
}

func (suite *ParsedArgsTestSuite) TestChecksum() {
	args := argparse.ParsedArgs{
		Parameters: map[string][]string{
			"v": {"k"},
			"k": {"kk"},
		},
		Flags: map[string]bool{
			"cleanup": true,
			"welcome": false,
		},
		Positional: []string{"1", "2", "3 4 5"},
	}

	suite.Equal(
		"p8bRasque_49ECe4B8wAml1cVeDj5pxomI0F94kl_b4",
		args.Checksum())
}

func (suite *ParsedArgsTestSuite) TestIsPositionalTime() {
	suite.T().Run("empty", func(t *testing.T) {
		args := argparse.ParsedArgs{
			Positional: []string{},
		}
		assert.False(t, args.IsPositionalTime())
	})

	suite.T().Run("full", func(t *testing.T) {
		args := argparse.ParsedArgs{
			Positional: []string{"1"},
		}
		assert.True(t, args.IsPositionalTime())
	})
}

func (suite *ParsedArgsTestSuite) TestCannotFindRequiredParameter() {
	args := argparse.ParsedArgs{
		Parameters: map[string][]string{
			"int1": {"1"},
		},
		Flags: map[string]bool{
			"flag1": true,
		},
	}

	suite.ErrorContains(
		args.Validate(suite.Context(), suite.flags, suite.params),
		"mandatory parameter")
}

func (suite *ParsedArgsTestSuite) TestCannotFindRequiredFlag() {
	args := argparse.ParsedArgs{
		Parameters: map[string][]string{
			"int1":  {"1"},
			"json1": {"[]"},
		},
		Flags: map[string]bool{
			"flag2": true,
		},
	}

	suite.ErrorContains(
		args.Validate(suite.Context(), suite.flags, suite.params),
		"mandatory flag")
}

func (suite *ParsedArgsTestSuite) TestUnknownParameter() {
	args := argparse.ParsedArgs{
		Parameters: map[string][]string{
			"int1":  {"1"},
			"json1": {"[]"},
			"x":     {"y"},
		},
		Flags: map[string]bool{
			"flag1": false,
			"flag2": true,
		},
	}

	suite.ErrorContains(
		args.Validate(suite.Context(), suite.flags, suite.params),
		"unknown parameter")
}

func (suite *ParsedArgsTestSuite) TestUnknownFlag() {
	args := argparse.ParsedArgs{
		Parameters: map[string][]string{
			"int1":  {"1"},
			"json1": {"[]"},
		},
		Flags: map[string]bool{
			"flag1": false,
			"flag2": true,
			"x":     true,
		},
	}

	suite.ErrorContains(
		args.Validate(suite.Context(), suite.flags, suite.params),
		"unknown flag")
}

func (suite *ParsedArgsTestSuite) TestIncorrectParameter() {
	args := argparse.ParsedArgs{
		Parameters: map[string][]string{
			"int1":  {"1"},
			"json1": {"{"},
		},
		Flags: map[string]bool{
			"flag1": false,
			"flag2": true,
		},
	}

	suite.ErrorContains(
		args.Validate(suite.Context(), suite.flags, suite.params),
		"invalid value for parameter")
}

func (suite *ParsedArgsTestSuite) TestAllParametersIncorrect() {
	args := argparse.ParsedArgs{
		Parameters: map[string][]string{
			"int1":  {"xxx"},
			"json1": {"{"},
		},
		Flags: map[string]bool{
			"flag1": false,
			"flag2": true,
		},
	}

	suite.ErrorContains(
		args.Validate(suite.Context(), suite.flags, suite.params),
		"invalid value for parameter")
}

func (suite *ParsedArgsTestSuite) TestValidateListOk() {
	args := argparse.ParsedArgs{
		Parameters: map[string][]string{
			"int1":  {"1", "2", "3"},
			"json1": {"{}"},
		},
		Flags: map[string]bool{
			"flag1": false,
			"flag2": true,
		},
	}

	suite.NoError(args.Validate(suite.Context(), suite.flags, suite.params))
}

func (suite *ParsedArgsTestSuite) TestValidateListFail() {
	args := argparse.ParsedArgs{
		Parameters: map[string][]string{
			"int1":  {"1", "xxx", "3"},
			"json1": {"{}"},
		},
		Flags: map[string]bool{
			"flag1": false,
			"flag2": true,
		},
	}

	err := args.Validate(suite.Context(), suite.flags, suite.params)

	suite.ErrorContains(err, "invalid value for parameter")
	suite.ErrorContains(err, "int1")
	suite.ErrorContains(err, "xxx")
}

func (suite *ParsedArgsTestSuite) TestGetParameterList() {
	testTable := map[string]string{
		"int1":  "1\nxxx yyy\n3",
		"json1": "{}",
		"":      "",
	}
	args := argparse.ParsedArgs{
		Parameters: map[string][]string{
			"int1":  {"1", "xxx yyy", "3"},
			"json1": {"{}"},
		},
		Flags: map[string]bool{
			"flag1": false,
			"flag2": true,
		},
	}

	for testValue, expected := range testTable {
		testValue := testValue
		expected := expected

		suite.T().Run(testValue, func(t *testing.T) {
			assert.Equal(t, expected, args.GetParameterList(testValue))
		})
	}
}

func (suite *ParsedArgsTestSuite) TestGetParameter() {
	testTable := map[string]string{
		"int1":  "3",
		"json1": "{}",
		"":      "",
	}
	args := argparse.ParsedArgs{
		Parameters: map[string][]string{
			"int1":  {"1", "xxx yyy", "3"},
			"json1": {"{}"},
		},
		Flags: map[string]bool{
			"flag1": false,
			"flag2": true,
		},
	}

	for testValue, expected := range testTable {
		testValue := testValue
		expected := expected

		suite.T().Run(testValue, func(t *testing.T) {
			assert.Equal(t, expected, args.GetParameter(testValue))
		})
	}
}

func (suite *ParsedArgsTestSuite) TestToSelfStringChunks() {
	args := argparse.ParsedArgs{
		Parameters: map[string][]string{
			"int1":  {"1", "xxx yyy", "3"},
			"json1": {"{}"},
		},
		Flags: map[string]bool{
			"flag1": false,
			"flag2": true,
		},
		Positional: []string{"p1", "p2"},
	}

	suite.Equal(
		[]string{"+flag2", "int1=1", "int1=xxx yyy", "int1=3", "json1={}"},
		args.ToSelfStringChunks())
}

func (suite *ParsedArgsTestSuite) TestToSlugString() {
	args := argparse.ParsedArgs{
		Parameters: map[string][]string{
			"int1":  {"1", "xxx yyy", "3"},
			"json1": {"{}"},
		},
		Flags: map[string]bool{
			"flag1": false,
			"flag2": true,
		},
		Positional: []string{"p1", "p2"},
	}

	suite.Equal(
		`0p1-0p2-2json1_-2int1_xxx-yyy-2int1_3-2int1_1-1flag2`,
		args.ToSlugString())
}

func TestParsedArgs(t *testing.T) {
	suite.Run(t, &ParsedArgsTestSuite{})
}
