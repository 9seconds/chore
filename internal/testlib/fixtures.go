package testlib

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	SaveSnapshotEnvVar     = "SAVE_SNAPSHOT"
	FixturesFilePermission = 0o644
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

func (suite *FixturesTestSuite) ReadPath(path string) []byte {
	suite.t.Helper()

	data, err := os.ReadFile(suite.FixturePath(path))
	require.NoError(suite.t, err)

	return data
}

func (suite *FixturesTestSuite) EnsureSnapshot(data []byte, path string) {
	suite.t.Helper()

	if _, ok := os.LookupEnv(SaveSnapshotEnvVar); ok {
		require.NoError(suite.t, os.WriteFile(
			filepath.Join("testdata", path),
			data,
			FixturesFilePermission,
		))
	}
}
