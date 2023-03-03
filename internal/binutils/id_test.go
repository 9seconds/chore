package binutils_test

import (
	"testing"

	"github.com/9seconds/chore/internal/binutils"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	sequence := make([]string, 1000)
	cardinality := make(map[string]bool)

	for i := 0; i < 1000; i++ {
		id := binutils.NewID()
		sequence[i] = id
		cardinality[id] = true
	}

	assert.IsIncreasing(t, sequence)
	assert.Len(t, cardinality, 1000)
}
