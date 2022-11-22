package filelock_test

import (
	"context"
	"testing"
	"time"

	"github.com/9seconds/chore/internal/filelock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SharedLockTestSuite struct {
	BaseLockTestSuite

	rlock1 filelock.Lock
	rlock2 filelock.Lock
	xlock  filelock.Lock
}

func (suite *SharedLockTestSuite) SetupTest() {
	suite.BaseLockTestSuite.SetupTest()

	var err error

	suite.rlock1, err = filelock.New(filelock.LockTypeShared, suite.path)
	require.NoError(suite.T(), err)

	suite.rlock2, err = filelock.New(filelock.LockTypeShared, suite.path)
	require.NoError(suite.T(), err)

	suite.xlock, err = filelock.New(filelock.LockTypeExclusive, suite.path)
	require.NoError(suite.T(), err)
}

func (suite *SharedLockTestSuite) TestAcquireRace() {
	suite.NoError(suite.rlock1.Lock(suite.Context()))
	suite.NoError(suite.rlock2.Lock(suite.Context()))

	failedCtx, cancel := context.WithTimeout(suite.Context(), 100*time.Millisecond)
	defer cancel()

	suite.ErrorContains(
		suite.xlock.Lock(failedCtx),
		"context deadline exceeded")

	suite.NoError(suite.rlock1.Unlock())
	suite.NoError(suite.rlock2.Unlock())
}

func (suite *SharedLockTestSuite) TestCannotAcquireIfSharedIsTaken() {
	suite.NoError(suite.xlock.Lock(suite.Context()))

	failedCtx, cancel := context.WithTimeout(suite.Context(), 100*time.Millisecond)
	defer cancel()

	suite.ErrorContains(
		suite.rlock1.Lock(failedCtx),
		"context deadline exceeded")
}

func (suite *SharedLockTestSuite) TestReentrantLock() {
	suite.NoError(suite.rlock1.Lock(suite.Context()))
	suite.NoError(suite.rlock1.Lock(suite.Context()))
	suite.NoError(suite.rlock1.Unlock())
	suite.NoError(suite.rlock1.Unlock())

	suite.NoError(suite.xlock.Lock(suite.Context()))
	suite.NoError(suite.xlock.Unlock())
}

func TestSharedLock(t *testing.T) {
	suite.Run(t, &SharedLockTestSuite{})
}
