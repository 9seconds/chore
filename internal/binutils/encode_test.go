package binutils_test

import (
	"testing"

	"github.com/9seconds/chore/internal/binutils"
	"github.com/stretchr/testify/assert"
)

func TestChain(t *testing.T) {
	assert.Equal(
		t,
		"1xyZoW4DhyYT9kuvZ297x_v1_AzJ6xby8KyeWfTRlwo",
		binutils.Chain("xx", "a", "b", "cd"))
}

func TestToString(t *testing.T) {
	assert.Equal(t, "AQID", binutils.ToString([]byte{1, 2, 3}))
}
