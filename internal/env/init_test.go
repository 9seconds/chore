package env_test

import (
	"context"
	"strings"
	"sync"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type EnvBaseTestSuite struct {
	suite.Suite

	wg        *sync.WaitGroup
	ctx       context.Context
	ctxCancel context.CancelFunc
	values    chan string
}

func (suite *EnvBaseTestSuite) SetupTest() {
	ctx, cancel := context.WithCancel(context.Background())

	suite.ctx = ctx
	suite.ctxCancel = cancel
	suite.values = make(chan string, 1)
	suite.wg = &sync.WaitGroup{}
}

func (suite *EnvBaseTestSuite) TearDownTest() {
	suite.ctxCancel()
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
