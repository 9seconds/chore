package git_test

import (
	"testing"

	"github.com/9seconds/chore/internal/git"
	"github.com/stretchr/testify/assert"
)

func TestAccessMode(t *testing.T) {
	testTable := map[string]bool{
		"no":           true,
		"if_undefined": true,
		"always":       true,
		"xx":           false,
	}

	for testValue, isValid := range testTable {
		testValue := testValue
		isValid := isValid

		t.Run(testValue, func(t *testing.T) {
			value := git.AccessMode(testValue)
			mode, err := git.GetAccessMode(testValue)

			assert.Equal(t, testValue, value.String())

			if isValid {
				assert.NoError(t, err)
				assert.True(t, value.Valid())
				assert.Equal(t, mode, value)
			} else {
				assert.Error(t, err)
				assert.False(t, value.Valid())
			}
		})
	}

	t.Run("default", func(t *testing.T) {
		mode, err := git.GetAccessMode("")

		assert.NoError(t, err)
		assert.Equal(t, git.AccessModeNo, mode)
	})
}
