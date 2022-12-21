package argparse_test

import (
	"encoding/hex"
	"testing"

	"github.com/9seconds/chore/internal/argparse"
	"github.com/stretchr/testify/assert"
)

func TestParsedArgs(t *testing.T) {
	args := argparse.ParsedArgs{
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

	assert.Equal(
		t,
		"5fd03abc5a84a11cf3a5ed4d1ed78e249ac5313053dd5536ab8f79cc21f8bdd3",
		hex.EncodeToString(args.Checksum()))
}
