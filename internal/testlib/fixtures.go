package testlib

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

type FixturesTestSuite struct {
	t *testing.T
}

func (suite *FixturesTestSuite) Setup(t *testing.T) {
	t.Helper()

	suite.t = t
}

func (suite *FixturesTestSuite) FixturePath(path string) string {
	suite.t.Helper()

	path = filepath.Join("testdata", path)

	require.FileExists(suite.t, path)

	return path
}
