package filelock_test

import (
	"context"
	"testing"
	"time"

	"github.com/9seconds/chore/internal/filelock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ExclusiveLockTestSuite struct {
	BaseLockTestSuite

	lock1 filelock.Lock
	lock2 filelock.Lock
}

func (suite *ExclusiveLockTestSuite) SetupTest() {
	suite.BaseLockTestSuite.SetupTest()

	var err error

	suite.lock1, err = filelock.New(filelock.LockTypeExclusive, suite.path)
	require.NoError(suite.T(), err)

	suite.lock2, err = filelock.New(filelock.LockTypeExclusive, suite.path)
	require.NoError(suite.T(), err)
}

func (suite *ExclusiveLockTestSuite) TestAcquireRace() {
	suite.NoError(suite.lock1.Lock(suite.Context()))

	failedCtx, cancel := context.WithTimeout(suite.Context(), 100*time.Millisecond)
	defer cancel()

	suite.ErrorContains(
		suite.lock2.Lock(failedCtx),
		"context deadline exceeded")

	suite.NoError(suite.lock1.Unlock())
}

func (suite *ExclusiveLockTestSuite) TestReleaseLockAndReaquire() {
	suite.NoError(suite.lock1.Lock(suite.Context()))

	go func() {
		<-time.After(50 * time.Millisecond)
		suite.NoError(suite.lock1.Unlock())
	}()

	suite.NoError(suite.lock2.Lock(suite.Context()))
	suite.NoError(suite.lock2.Unlock())
}

func (suite *ExclusiveLockTestSuite) TestReentrantLock() {
	ctx, cancel := context.WithTimeout(suite.Context(), time.Second)
	defer cancel()

	suite.NoError(suite.lock1.Lock(ctx))
	suite.NoError(suite.lock1.Lock(ctx))
	suite.NoError(suite.lock1.Unlock())
	suite.NoError(suite.lock2.Unlock())
}

func TestExclusiveLock(t *testing.T) {
	suite.Run(t, &ExclusiveLockTestSuite{})
}
