package env_test

import (
	"os"
	"testing"

	"github.com/9seconds/chore/internal/env"
	"github.com/Showmax/go-fqdn"
	"github.com/stretchr/testify/suite"
)

type GenerateHostnameTestSuite struct {
	EnvBaseTestSuite
}

func (suite *GenerateHostnameTestSuite) TestHostname() {
	value, err := os.Hostname()
	if err != nil {
		suite.T().Skipf("Hostname is not available: %v", err)
	}

	env.GenerateHostname(suite.Context(), suite.values, suite.wg)
	data := suite.Collect()

	suite.Equal(value, data[env.EnvHostname])
}

func (suite *GenerateHostnameTestSuite) TestFQDNHostname() {
	value, err := fqdn.FqdnHostname()
	if err != nil {
		suite.T().Skipf("FQDN Hostname is not available: %v", err)
	}

	env.GenerateHostname(suite.Context(), suite.values, suite.wg)
	data := suite.Collect()

	suite.Equal(value, data[env.EnvHostnameFQDN])
}

func (suite *GenerateHostnameTestSuite) TestWithEnv() {
	suite.Setenv(env.EnvHostname, "xx")
	suite.Setenv(env.EnvHostnameFQDN, "yy")

	env.GenerateHostname(suite.Context(), suite.values, suite.wg)
	data := suite.Collect()

	suite.Empty(data)
}

func TestGenerateHostname(t *testing.T) {
	suite.Run(t, &GenerateHostnameTestSuite{})
}
