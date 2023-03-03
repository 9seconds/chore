package paths_test

import (
	"testing"

	"github.com/9seconds/chore/internal/paths"
	"github.com/stretchr/testify/assert"
)

func TestTemp(t *testing.T) {
	dir1, err := paths.TempDir()
	assert.NoError(t, err)
	assert.DirExists(t, dir1)

	dir2, err := paths.TempDir()
	assert.NoError(t, err)
	assert.DirExists(t, dir2)

	dir3, err := paths.TempDir()
	assert.NoError(t, err)
	assert.DirExists(t, dir3)

	paths.TempDirCleanup()
	assert.NoDirExists(t, dir1)
	assert.NoDirExists(t, dir2)
	assert.NoDirExists(t, dir3)
}
