package cli_test

import (
	"testing"

	"github.com/9seconds/chore/internal/cli"
	"github.com/9seconds/chore/internal/env"
	"github.com/stretchr/testify/suite"
)

type NamespaceTestSuite struct {
	suite.Suite

	param cli.Namespace
}

func (suite *NamespaceTestSuite) TestParse() {
	suite.NoError(suite.param.UnmarshalText([]byte("123")))
	suite.Equal("123", suite.param.Value())
}

func (suite *NamespaceTestSuite) TestParseMagicValueNoEnv() {
	suite.ErrorContains(
		suite.param.UnmarshalText([]byte(cli.MagicValue)),
		"namespace is magic")
}

func (suite *NamespaceTestSuite) TestParseMagicValueWithEnv() {
	suite.T().Setenv(env.EnvNamespace, "lalala")

	suite.NoError(suite.param.UnmarshalText([]byte(cli.MagicValue)))
	suite.Equal("lalala", suite.param.Value())
}

func TestNamespace(t *testing.T) {
	suite.Run(t, &NamespaceTestSuite{})
}
