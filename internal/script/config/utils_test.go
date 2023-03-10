package config_test

import (
	"testing"

	"github.com/9seconds/chore/internal/script/config"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeName(t *testing.T) {
	testTable := map[string]string{
		"":       "",
		"xx":     "xx",
		"X":      "x",
		"A B C":  "a_b_c",
		"a-B--c": "a_b__c",
	}

	for testName, expected := range testTable {
		testName := testName
		expected := expected

		t.Run(testName, func(t *testing.T) {
			assert.Equal(t, expected, config.NormalizeName(testName))
		})
	}
}
