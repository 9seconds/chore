package git

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemoize(t *testing.T) {
	callback := memoize(func() (time.Time, error) {
		return time.Now(), nil
	})

	value1, err := callback()
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	value2, err := callback()
	assert.NoError(t, err)

	assert.Equal(t, value1, value2)
}
