package env_test

import (
	"runtime"
	"testing"

	"github.com/9seconds/chore/internal/env"
	"github.com/stretchr/testify/suite"
)

type GenerateOSTestSuite struct {
	EnvBaseTestSuite
}

func (suite *GenerateOSTestSuite) TestNoEnv() {
	env.GenerateOS(suite.Context(), suite.values, suite.wg)
	data := suite.Collect()

	suite.Equal(runtime.GOOS, data[env.EnvOSType])
	suite.Equal(runtime.GOARCH, data[env.EnvOSArch])
}

func TestGenerateOS(t *testing.T) {
	t.Parallel()
	suite.Run(t, &GenerateOSTestSuite{})
}
