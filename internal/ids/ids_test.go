package ids_test

import (
	"testing"

	"github.com/9seconds/chore/internal/ids"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	sequence := make([]string, 1000)
	cardinality := make(map[string]bool)

	for i := 0; i < 1000; i++ {
		id := ids.New()
		sequence[i] = id
		cardinality[id] = true
	}

	assert.IsIncreasing(t, sequence)
	assert.Len(t, cardinality, 1000)
}

func TestChain(t *testing.T) {
	assert.Equal(
		t,
		"GZEWAmDCaKTvpTtqdvTfkq_gWqS4KcsHKwkiMEpAOlI",
		ids.Chain("xx", "a", "b", "cd"))
}

func TestEncode(t *testing.T) {
	assert.Equal(t, "AQID", ids.Encode([]byte{1, 2, 3}))
}
