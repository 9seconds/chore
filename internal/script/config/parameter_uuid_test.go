package config_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/9seconds/chore/internal/script/config"
	"github.com/9seconds/chore/internal/testlib"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ParameterUUIDTestSuite struct {
	suite.Suite

	testlib.CtxTestSuite
}

func (suite *ParameterUUIDTestSuite) SetupTest() {
	suite.CtxTestSuite.Setup(suite.T())
}

func (suite *ParameterUUIDTestSuite) TestRequired() {
	testTable := []bool{true, false}

	for _, testValue := range testTable {
		testValue := testValue

		suite.T().Run(strconv.FormatBool(testValue), func(t *testing.T) {
			param, err := config.NewUUID("", testValue, nil)
			assert.NoError(t, err)
			assert.Equal(t, testValue, param.Required())
		})
	}
}

func (suite *ParameterUUIDTestSuite) TestType() {
	param, err := config.NewUUID("", false, nil)
	suite.NoError(err)
	suite.Equal(config.ParameterUUID, param.Type())
}

func (suite *ParameterUUIDTestSuite) TestIncorrectVersion() {
	testTable := map[string]bool{
		"":      false,
		"{":     false,
		"-1":    false,
		"10":    false,
		"255":   false,
		"11111": false,

		"1": true,
		"2": false,
		"3": true,
		"4": true,
		"5": true,
		"6": true,
		"7": true,
	}

	for testName, isValid := range testTable {
		testName := testName
		isValid := isValid

		suite.T().Run(testName, func(t *testing.T) {
			_, err := config.NewUUID("", false, map[string]string{
				"version": testName,
			})

			if isValid {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, "incorrect version")
			}
		})
	}
}

func (suite *ParameterUUIDTestSuite) TestValidation() {
	wrongValues := []string{
		"",
		"{",
		"-1",
		"10",
		"255",
		"11111",
	}

	uuids := []uuid.UUID{
		uuid.Must(uuid.NewV1()),
		uuid.NewV3(uuid.NamespaceDNS, "xx"),
		uuid.Must(uuid.NewV4()),
		uuid.NewV5(uuid.NamespaceDNS, "xx"),
		uuid.Must(uuid.NewV6()),
		uuid.Must(uuid.NewV7()),
	}

	v4Only, err := config.NewUUID("", false, map[string]string{
		"version": "4",
	})
	suite.NoError(err)

	allOk, err := config.NewUUID("", false, nil)
	suite.NoError(err)

	for _, testValue := range wrongValues {
		testValue := testValue

		suite.T().Run(testValue, func(t *testing.T) {
			assert.ErrorContains(
				t,
				allOk.Validate(suite.Context(), testValue),
				"cannot parse uuid")
		})
	}

	for _, uuidValue := range uuids {
		version := uuidValue.Version()
		testValue := uuidValue.String()

		suite.T().Run(fmt.Sprintf("v4 - %d", version), func(t *testing.T) {
			err := v4Only.Validate(suite.Context(), testValue)

			if version == 4 {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, "incorrect uuid version")
			}
		})

		suite.T().Run(fmt.Sprintf("all - %d", version), func(t *testing.T) {
			assert.NoError(t, allOk.Validate(suite.Context(), testValue))
		})
	}
}

func TestParameterUUID(t *testing.T) {
	suite.Run(t, &ParameterUUIDTestSuite{})
}
