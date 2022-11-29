package cli_test

import (
	"testing"
	"time"

	"github.com/9seconds/chore/internal/cli"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ContextTestSuite struct {
	suite.Suite
	testlib.CtxTestSuite

	ctx  cli.Context
	lock *testlib.FileLockMock
}

func (suite *ContextTestSuite) SetupTest() {
	suite.CtxTestSuite.Setup(suite.T())

	suite.lock = &testlib.FileLockMock{}
	suite.ctx = cli.NewContext(suite.Context()).WithLock(suite.lock)

	suite.T().Cleanup(func() {
		suite.ctx.Close()

		suite.lock.AssertExpectations(suite.T())
	})
}

func (suite *ContextTestSuite) TestCloseContextClosesContext() {
	suite.ctx = cli.NewContext(suite.Context())
	suite.NoError(suite.ctx.Close())

	select {
	case <-suite.ctx.Done():
	default:
		suite.Fail("context is not closed")
	}
}

func (suite *ContextTestSuite) TestWithTimeout() {
	suite.lock.On("Unlock").Return(nil)

	ctx := suite.ctx.WithTimeout(10 * time.Millisecond)

	time.Sleep(50 * time.Millisecond)

	select {
	case <-ctx.Done():
	default:
		suite.Fail("context is not closed")
	}
}

func (suite *ContextTestSuite) TestWithLock() {
	suite.lock.
		On("Unlock").
		Return(nil)
	suite.lock.
		On("Lock", mock.Anything).
		Return(nil)

	suite.NoError(suite.ctx.Start())
}

func TestContext(t *testing.T) {
	suite.Run(t, &ContextTestSuite{})
}
