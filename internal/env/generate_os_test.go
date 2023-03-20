package env_test

import (
	"runtime"
	"testing"

	"github.com/9seconds/chore/internal/env"
	"github.com/stretchr/testify/suite"
)

type GenerateOSTestSuite struct {
	BaseTestSuite
}

func (suite *GenerateOSTestSuite) TestNo() {
	env.GenerateOS(suite.Context(), suite.values, suite.wg)
	data := suite.Collect()

	suite.Equal(runtime.GOOS, data[env.OSType])
	suite.Equal(runtime.GOARCH, data[env.OSArch])
}

func TestGenerateOS(t *testing.T) {
	suite.Run(t, &GenerateOSTestSuite{})
}
