package argparse_test

import (
	"encoding/hex"
	"testing"

	"github.com/9seconds/chore/internal/argparse"
	"github.com/stretchr/testify/suite"
)

type ParsedArgsTestSuite struct {
	suite.Suite
}

func (suite *ParsedArgsTestSuite) TestChecksum() {
	args := argparse.ParsedArgs{
		Keywords: map[string]string{
			"v": "k",
			"k": "2",
		},
		Positional: []string{"1", "2", "3 4 5"},
	}

	suite.Equal(
		"c773754b585fabf095af52a3113b0e24661fe60e17719cd12b217284a67dce7a",
		hex.EncodeToString(args.Checksum()))
}

func TestParsedArgs(t *testing.T) {
	t.Parallel()
	suite.Run(t, &ParsedArgsTestSuite{})
}
