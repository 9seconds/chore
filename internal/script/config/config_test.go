package config_test

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/9seconds/chore/internal/script/config"
	"github.com/stretchr/testify/assert"
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
			buf := strings.NewReader(fmt.Sprintf("network = %t", testValue))

			conf, err := config.Parse(buf)
			assert.NoError(t, err)
			assert.Equal(t, testValue, conf.Network)
		})
	}
}

func (suite *ConfigTestSuite) TestParseDescription() {
	buf := strings.NewReader("description = 'xxy'")

	conf, err := config.Parse(buf)
	suite.NoError(err)
	suite.Equal("xxy", conf.Description)
}

func (suite *ConfigTestSuite) TestParameter() {
	tableTest := []string{
		config.ParameterInteger,
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
		config.ParameterGit,
	}

	for _, testValue := range tableTest {
		testValue := testValue

		suite.T().Run(testValue, func(t *testing.T) {
			specRaw := `
[parameters.param]
type = "%s"
required = true
spec = %s`
			spec := "{}"

			if testValue == config.ParameterEnum {
				spec = `{ choices = "xx" }`
			}

			buf := strings.NewReader(fmt.Sprintf(specRaw, testValue, spec))

			conf, err := config.Parse(buf)
			suite.NoError(err)
			suite.Len(conf.Parameters, 1)
			suite.True(conf.Parameters["param"].Required())
			suite.Equal(testValue, conf.Parameters["param"].Type())
		})
	}
}

func (suite *ConfigTestSuite) TestUnknownParameterType() {
	configRaw := `
[parameters.param]
type = "xxx"`
	buf := strings.NewReader(configRaw)

	_, err := config.Parse(buf)
	suite.ErrorContains(err, "unknown parameter type")
}

func (suite *ConfigTestSuite) TestCannotInitializeParameter() {
	configRaw := `
[parameters.param]
type = "integer"
required = true
spec = { max = "x" }`
	buf := strings.NewReader(configRaw)

	_, err := config.Parse(buf)
	suite.ErrorContains(err, "cannot initialize parameter")
}

func (suite *ConfigTestSuite) TestIncorrectJSON() {
	buf := bytes.NewBuffer([]byte("x"))

	_, err := config.Parse(buf)
	suite.ErrorContains(err, "cannot parse TOML config")
}

func (suite *ConfigTestSuite) TestBrokenReader() {
	_, err := config.Parse(iotest.ErrReader(io.ErrUnexpectedEOF))
	suite.ErrorIs(err, io.ErrUnexpectedEOF)
}

func TestConfig(t *testing.T) {
	suite.Run(t, &ConfigTestSuite{})
}
