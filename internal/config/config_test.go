package config_test

import (
	"sort"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/9seconds/chore/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

func (suite *ConfigTestSuite) TestErrorRead() {
	reader := iotest.TimeoutReader(strings.NewReader("[vault]"))
	_, err := config.ReadConfig(reader)

	suite.ErrorContains(err, "cannot parse TOML config")
	suite.ErrorIs(err, iotest.ErrTimeout)
}

func (suite *ConfigTestSuite) TestBadPasswords() {
	reader := strings.NewReader(`
		[vault.xx]
		y = 1
	`)
	_, err := config.ReadConfig(reader)

	suite.ErrorContains(err, "cannot parse TOML config")
}

func (suite *ConfigTestSuite) TestBadEnv() {
	reader := strings.NewReader(`
		[env.x]
		t = 1
	`)
	_, err := config.ReadConfig(reader)

	suite.ErrorContains(err, "cannot parse TOML config")
}

func (suite *ConfigTestSuite) TestOk() {
	reader := strings.NewReader(`
		[env.x]
		y = "1"

		[vault]
		y = "1"
		z = "2"
	`)

	conf, err := config.ReadConfig(reader)
	suite.NoError(err)
	suite.Equal(conf.Vault, map[string]string{"y": "1", "z": "2"})
	suite.Equal(conf.Env, map[string]map[string]string{"x": {"y": "1"}})
}

func (suite *ConfigTestSuite) TestEnviron() {
	reader := strings.NewReader(`
		[env.x]
		y = "1"
		key = "2"
	`)

	conf, err := config.ReadConfig(reader)
	suite.NoError(err)

	testTable := map[string][]string{
		"z": {},
		"x": {"y=1", "key=2"},
	}

	for testValue, expected := range testTable {
		testValue := testValue
		expected := expected

		suite.T().Run(testValue, func(t *testing.T) {
			actual := conf.Environ(testValue)

			sort.Strings(actual)
			sort.Strings(expected)

			assert.Equal(t, expected, actual)
		})
	}
}

func TestConfig(t *testing.T) {
	suite.Run(t, &ConfigTestSuite{})
}
