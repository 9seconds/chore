package filelock_test

import (
	"testing"

	"github.com/9seconds/chore/internal/filelock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type NoopLockTestSuite struct {
	BaseLockTestSuite

	lock filelock.Lock
}

func (suite *NoopLockTestSuite) SetupTest() {
	lock, err := filelock.New(filelock.LockTypeNo, suite.path)
	require.NoError(suite.T(), err)

	suite.lock = lock
}

func (suite *NoopLockTestSuite) TestAcquire() {
	suite.NoError(suite.lock.Lock(suite.Context()))
	suite.NoError(suite.lock.Unlock())
}

func TestNoopLock(t *testing.T) {
	suite.Run(t, &NoopLockTestSuite{})
}
