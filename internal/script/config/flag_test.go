package config_test

import (
	"strconv"
	"testing"

	"github.com/9seconds/chore/internal/script/config"
	"github.com/stretchr/testify/assert"
)

func TestFlag(t *testing.T) {
	for _, v := range []bool{true, false} {
		v := v

		t.Run(strconv.FormatBool(v), func(t *testing.T) {
			param := config.NewFlag("lalala", v)

			if v {
				assert.True(t, param.Required())
			} else {
				assert.False(t, param.Required())
			}

			assert.Contains(t, param.String(), "lalala")
		})
	}
}
