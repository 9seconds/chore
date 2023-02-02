package git

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemoize(t *testing.T) {
	fn := memoize(func() (time.Time, error) {
		return time.Now(), nil
	})

	value1, err := fn()
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	value2, err := fn()
	assert.NoError(t, err)

	assert.Equal(t, value1, value2)
}
