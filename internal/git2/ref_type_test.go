package git2_test

import (
	"testing"

	"github.com/9seconds/chore/internal/git2"
	"github.com/stretchr/testify/assert"
)

func TestRefType(t *testing.T) {
	testTable := map[string]bool{
		"revision": true,
		"tag":      true,
		"branch":   true,
		"remote":   true,
		"note":     true,
		"commit":   true,
		"xx":       false,
	}

	for testValue, isValid := range testTable {
		testValue := testValue
		isValid := isValid

		t.Run(testValue, func(t *testing.T) {
			value := git2.RefType(testValue)
			ref, err := git2.GetRefType(testValue)

			assert.Equal(t, testValue, value.String())

			if isValid {
				assert.NoError(t, err)
				assert.True(t, value.Valid())
				assert.Equal(t, ref, value)
			} else {
				assert.Error(t, err)
				assert.False(t, value.Valid())
			}
		})
	}

	t.Run("default", func(t *testing.T) {
		ref, err := git2.GetRefType("")

		assert.NoError(t, err)
		assert.Equal(t, git2.RefTypeRevision, ref)
	})
}
