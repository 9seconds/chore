package config_test

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"testing"
	"testing/iotest"

	"github.com/9seconds/chore/internal/config"
	"github.com/alecthomas/assert/v2"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

func (suite *ConfigTestSuite) TestParseNetwork() {
	testTable := []bool{true, false}

	for _, testValue := range testTable {
		testValue := testValue

		suite.T().Run(strconv.FormatBool(testValue), func(t *testing.T) {
			configRaw := map[string]bool{
				"network": testValue,
			}
			data, _ := json.Marshal(configRaw)
			buf := bytes.NewBuffer(data)

			conf, err := config.Parse(buf)
			assert.NoError(t, err)
			assert.Equal(t, testValue, conf.Network)
		})
	}
}

func (suite *ConfigTestSuite) TestParseDescription() {
	configRaw := map[string]string{
		"description": "xxy",
	}
	data, _ := json.Marshal(configRaw) //nolint: errchkjson
	buf := bytes.NewBuffer(data)

	conf, err := config.Parse(buf)
	suite.NoError(err)
	suite.Equal("xxy", conf.Description)
}

func (suite *ConfigTestSuite) TestParameter() {
	tableTest := []string{
		config.ParameterInteger,
		config.ParameterBool,
		config.ParameterString,
		config.ParameterFloat,
		config.ParameterURL,
		config.ParameterEmail,
		config.ParameterEnum,
		config.ParameterBase64,
		config.ParameterHex,
		config.ParameterHostname,
		config.ParameterMac,
		config.ParameterJSON,
		config.ParameterXML,
		config.ParameterUUID,
		config.ParameterDirectory,
		config.ParameterFile,
		config.ParameterSemver,
		config.ParameterDatetime,
	}

	for _, testValue := range tableTest {
		testValue := testValue

		suite.T().Run(testValue, func(t *testing.T) {
			spec := map[string]interface{}{}

			if testValue == config.ParameterEnum {
				spec = map[string]interface{}{
					"choices": "xx",
				}
			}

			configRaw := map[string]interface{}{
				"parameters": map[string]interface{}{
					"param": map[string]interface{}{
						"type":     testValue,
						"required": true,
						"spec":     spec,
					},
				},
			}
			data, _ := json.Marshal(configRaw)
			buf := bytes.NewBuffer(data)

			conf, err := config.Parse(buf)
			suite.NoError(err)
			suite.Len(conf.Parameters, 1)
			suite.True(conf.Parameters["param"].Required())
			suite.Equal(testValue, conf.Parameters["param"].Type())
		})
	}
}

func (suite *ConfigTestSuite) TestUnknownParameterType() {
	configRaw := map[string]interface{}{
		"parameters": map[string]interface{}{
			"param": map[string]interface{}{
				"type":     "xxx",
				"required": true,
				"spec":     map[string]interface{}{},
			},
		},
	}
	data, _ := json.Marshal(configRaw) //nolint: errchkjson
	buf := bytes.NewBuffer(data)

	_, err := config.Parse(buf)
	suite.ErrorContains(err, "unknown parameter type")
}

func (suite *ConfigTestSuite) TestCannotInitializeParameter() {
	configRaw := map[string]interface{}{
		"parameters": map[string]interface{}{
			"param": map[string]interface{}{
				"type":     "integer",
				"required": true,
				"spec": map[string]interface{}{
					"max": "x",
				},
			},
		},
	}
	data, _ := json.Marshal(configRaw) //nolint: errchkjson
	buf := bytes.NewBuffer(data)

	_, err := config.Parse(buf)
	suite.ErrorContains(err, "cannot initialize parameter")
}

func (suite *ConfigTestSuite) TestIncorrectParameterName() {
	configRaw := map[string]interface{}{
		"parameters": map[string]interface{}{
			"param 11": map[string]interface{}{
				"type":     "integer",
				"required": true,
				"spec":     map[string]interface{}{},
			},
		},
	}
	data, _ := json.Marshal(configRaw) //nolint: errchkjson
	buf := bytes.NewBuffer(data)

	_, err := config.Parse(buf)
	suite.ErrorContains(err, "incorrect parameter name")
}

func (suite *ConfigTestSuite) TestIncorrectJSON() {
	buf := bytes.NewBuffer([]byte("x"))

	_, err := config.Parse(buf)
	suite.ErrorContains(err, "cannot parse JSON config")
}

func (suite *ConfigTestSuite) TestBrokenReader() {
	_, err := config.Parse(iotest.ErrReader(io.ErrUnexpectedEOF))
	suite.ErrorIs(err, io.ErrUnexpectedEOF)
}

func TestConfig(t *testing.T) {
	suite.Run(t, &ConfigTestSuite{})
}
