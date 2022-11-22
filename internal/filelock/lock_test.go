package filelock_test

import (
	"testing"

	"github.com/9seconds/chore/internal/filelock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type NewLockTestSuite struct {
	BaseLockTestSuite
}

func (suite *NewLockTestSuite) TestAbsentPath() {
	testTable := map[filelock.LockType]bool{
		filelock.LockTypeNo:        true,
		filelock.LockTypeExclusive: false,
		filelock.LockTypeShared:    false,
	}

	for testName, isValid := range testTable {
		testName := testName
		isValid := isValid

		suite.T().Run(testName.String(), func(t *testing.T) {
			_, err := filelock.New(testName, "xx")

			if isValid {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, "cannot stat path")
			}
		})
	}
}

func (suite *NewLockTestSuite) TestUnknownLockType() {
	_, err := filelock.New(100, suite.path)
	suite.ErrorContains(err, "unknown lock type")
}

func (suite *NewLockTestSuite) TestLockWithCorrectPath() {
	testTable := []filelock.LockType{
		filelock.LockTypeNo,
		filelock.LockTypeExclusive,
		filelock.LockTypeShared,
	}

	for _, testValue := range testTable {
		testValue := testValue

		suite.T().Run(testValue.String(), func(t *testing.T) {
			_, err := filelock.New(testValue, suite.path)
			assert.NoError(t, err)
		})
	}
}

func TestNewLock(t *testing.T) {
	suite.Run(t, &NewLockTestSuite{})
}
