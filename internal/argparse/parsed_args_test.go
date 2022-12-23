package argparse_test

import (
	"encoding/hex"
	"sort"
	"testing"

	"github.com/9seconds/chore/internal/argparse"
	"github.com/stretchr/testify/suite"
)

type ParsedArgsTestSuite struct {
	suite.Suite

	args argparse.ParsedArgs
}

func (suite *ParsedArgsTestSuite) SetupTest() {
	suite.args = argparse.ParsedArgs{
		Parameters: map[string]string{
			"v": "k",
			"k": "2",
		},
		Flags: map[string]bool{
			"cleanup": true,
			"welcome": false,
		},
		Positional: []string{"1", "2", "3 4 5"},
	}
}

func (suite *ParsedArgsTestSuite) TestChecksum() {
	suite.Equal(
		"5fd03abc5a84a11cf3a5ed4d1ed78e249ac5313053dd5536ab8f79cc21f8bdd3",
		hex.EncodeToString(suite.args.Checksum()))
}

func (suite *ParsedArgsTestSuite) TestOptions() {
	options := suite.args.Options()

	sort.Strings(options)

	suite.Equal(
		[]string{"+cleanup", "-welcome", "k=2", "v=k"},
		options)
}

func TestParsedArgs(t *testing.T) {
	suite.Run(t, &ParsedArgsTestSuite{})
}
