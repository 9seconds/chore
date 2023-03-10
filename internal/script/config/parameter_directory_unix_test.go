//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || zos

package config_test

import (
	"io/fs"
	"os"
	"testing"

	"github.com/9seconds/chore/internal/script/config"
	"github.com/stretchr/testify/assert"
)

const (
	permUnixX fs.FileMode = (1 << (1 + iota)) << 5
	permUnixW
	permUnixR
)

func (suite *ParameterDirectoryTestSuite) TestPermissionsOk() {
	testTable := []string{
		"readable",
		"writable",
		"executable",
	}

	for _, presentPermission := range testTable {
		presentPermission := presentPermission

		suite.T().Run(presentPermission, func(t *testing.T) {
			var fileMode fs.FileMode

			switch presentPermission {
			case "readable":
				fileMode = permUnixR
			case "writable":
				fileMode = permUnixW
			case "executable":
				fileMode = permUnixX
			}

			assert.NoError(t, os.Chmod(suite.dir, fileMode))

			param, err := config.NewDirectory("", false, map[string]string{
				presentPermission: "true",
			})
			assert.NoError(t, err)
			assert.NoError(t, param.Validate(suite.Context(), suite.dir))
		})
	}
}

func (suite *ParameterDirectoryTestSuite) TestPermissionsNOk() {
	testTable := []string{
		"readable",
		"writable",
		"executable",
	}

	for _, absentPermission := range testTable {
		absentPermission := absentPermission

		suite.T().Run(absentPermission, func(t *testing.T) {
			fileMode := permUnixR | permUnixX | permUnixW

			switch absentPermission {
			case "readable":
				fileMode &^= permUnixR
			case "writable":
				fileMode &^= permUnixW
			case "executable":
				fileMode &^= permUnixX
			}

			assert.NoError(t, os.Chmod(suite.dir, fileMode))

			param, err := config.NewDirectory("", false, map[string]string{
				absentPermission: "true",
			})
			assert.NoError(t, err)
			assert.ErrorContains(
				t,
				param.Validate(suite.Context(), suite.dir),
				"incorrect user permissions")
		})
	}
}
