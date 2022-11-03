package env_test

import (
	"strings"
	"sync"

	"github.com/9seconds/chore/internal/testlib"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type EnvBaseTestSuite struct {
	suite.Suite

	testlib.CtxTestSuite

	wg     *sync.WaitGroup
	values chan string
}

func (suite *EnvBaseTestSuite) SetupTest() {
	suite.CtxTestSuite.Setup(suite.T())

	suite.values = make(chan string, 1)
	suite.wg = &sync.WaitGroup{}
}

func (suite *EnvBaseTestSuite) TearDownTest() {
	suite.wg.Wait()
}

func (suite *EnvBaseTestSuite) Collect() map[string]string {
	go func() {
		suite.wg.Wait()
		close(suite.values)
	}()

	collected := make(map[string]string)

	for text := range suite.values {
		name, value, found := strings.Cut(text, "=")
		require.True(suite.T(), found)

		collected[name] = value
	}

	return collected
}

func (suite *EnvBaseTestSuite) Setenv(name, value string) {
	suite.T().Setenv(name, value)
}
