package argparse_test

import (
	"encoding/hex"
	"sort"
	"testing"

	"github.com/9seconds/chore/internal/argparse"
	"github.com/9seconds/chore/internal/config"
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
		Parameters: map[string]string{
			"v": "k",
			"k": "kk",
		},
		Flags: map[string]argparse.FlagValue{
			"cleanup": argparse.FlagTrue,
			"welcome": argparse.FlagFalse,
		},
		Positional: []string{"1", "2", "3 4 5"},
	}

	suite.Equal(
		"896771bd25f4b67e7a71eb593659d63e17405d645442905a1b9c8d5a2000be05",
		hex.EncodeToString(args.Checksum()))
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

func (suite *ParsedArgsTestSuite) TestOptionsAsCli() {
	args := argparse.ParsedArgs{
		Parameters: map[string]string{
			"v": "k",
			"k": "kk",
		},
		Flags: map[string]argparse.FlagValue{
			"cleanup": argparse.FlagTrue,
			"welcome": argparse.FlagFalse,
		},
		Positional: []string{"1", "2", "3 4 5"},
	}
	options := args.OptionsAsCli()

	sort.Strings(options)

	suite.Equal([]string{"+cleanup", "@welcome", "k=kk", "v=k"}, options)
}

func (suite *ParsedArgsTestSuite) TestCannotFindRequiredParameter() {
	args := argparse.ParsedArgs{
		Parameters: map[string]string{
			"int1": "1",
		},
		Flags: map[string]argparse.FlagValue{
			"flag1": argparse.FlagTrue,
		},
	}

	suite.ErrorContains(
		args.Validate(suite.Context(), suite.flags, suite.params),
		"mandatory parameter")
}

func (suite *ParsedArgsTestSuite) TestCannotFindRequiredFlag() {
	args := argparse.ParsedArgs{
		Parameters: map[string]string{
			"int1":  "1",
			"json1": "[]",
		},
		Flags: map[string]argparse.FlagValue{
			"flag2": argparse.FlagTrue,
		},
	}

	suite.ErrorContains(
		args.Validate(suite.Context(), suite.flags, suite.params),
		"mandatory flag")
}

func (suite *ParsedArgsTestSuite) TestUnknownParameter() {
	args := argparse.ParsedArgs{
		Parameters: map[string]string{
			"int1":  "1",
			"json1": "[]",
			"x":     "y",
		},
		Flags: map[string]argparse.FlagValue{
			"flag1": argparse.FlagFalse,
			"flag2": argparse.FlagTrue,
		},
	}

	suite.ErrorContains(
		args.Validate(suite.Context(), suite.flags, suite.params),
		"unknown parameter")
}

func (suite *ParsedArgsTestSuite) TestUnknownFlag() {
	args := argparse.ParsedArgs{
		Parameters: map[string]string{
			"int1":  "1",
			"json1": "[]",
		},
		Flags: map[string]argparse.FlagValue{
			"flag1": argparse.FlagFalse,
			"flag2": argparse.FlagTrue,
			"x":     argparse.FlagTrue,
		},
	}

	suite.ErrorContains(
		args.Validate(suite.Context(), suite.flags, suite.params),
		"unknown flag")
}

func (suite *ParsedArgsTestSuite) TestIncorrectParameter() {
	args := argparse.ParsedArgs{
		Parameters: map[string]string{
			"int1":  "1",
			"json1": "{",
		},
		Flags: map[string]argparse.FlagValue{
			"flag1": argparse.FlagFalse,
			"flag2": argparse.FlagTrue,
		},
	}

	suite.ErrorContains(
		args.Validate(suite.Context(), suite.flags, suite.params),
		"invalid value for parameter")
}

func (suite *ParsedArgsTestSuite) TestAllParametersIncorrect() {
	args := argparse.ParsedArgs{
		Parameters: map[string]string{
			"int1":  "xxx",
			"json1": "{",
		},
		Flags: map[string]argparse.FlagValue{
			"flag1": argparse.FlagFalse,
			"flag2": argparse.FlagTrue,
		},
	}

	suite.ErrorContains(
		args.Validate(suite.Context(), suite.flags, suite.params),
		"invalid value for parameter")
}

func (suite *ParsedArgsTestSuite) TestValidateListOk() {
	args := argparse.ParsedArgs{
		Parameters: map[string]string{
			"int1":  "1:2:3",
			"json1": "{}",
		},
		Flags: map[string]argparse.FlagValue{
			"flag1": argparse.FlagFalse,
			"flag2": argparse.FlagTrue,
		},
		ListDelimiter: argparse.DefaultListDelimiter,
	}

	suite.NoError(args.Validate(suite.Context(), suite.flags, suite.params))
}

func (suite *ParsedArgsTestSuite) TestValidateListFail() {
	args := argparse.ParsedArgs{
		Parameters: map[string]string{
			"int1":  "1:xxx:3",
			"json1": "{}",
		},
		Flags: map[string]argparse.FlagValue{
			"flag1": argparse.FlagFalse,
			"flag2": argparse.FlagTrue,
		},
		ListDelimiter: argparse.DefaultListDelimiter,
	}

	err := args.Validate(suite.Context(), suite.flags, suite.params)

	suite.ErrorContains(err, "invalid value for parameter")
	suite.ErrorContains(err, "int1")
	suite.ErrorContains(err, "xxx")
}

func TestParsedArgs(t *testing.T) {
	suite.Run(t, &ParsedArgsTestSuite{})
}
