package config_test

import (
	"strings"
	"testing"
	"testing/iotest"

	"github.com/9seconds/chore/internal/config"
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

func (suite *ConfigTestSuite) TestOk() {
	reader := strings.NewReader(`
		[vault]
		y = "1"
		z = "2"
	`)

	conf, err := config.ReadConfig(reader)
	suite.NoError(err)
	suite.Subset(conf.Vault, map[string]string{"y": "1", "z": "2"})
}

func TestConfig(t *testing.T) {
	suite.Run(t, &ConfigTestSuite{})
}
