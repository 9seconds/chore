package env_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/9seconds/chore/internal/env"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"
)

type GenerateNetworkBaseTestSuite struct {
	EnvBaseTestSuite
}

func (suite *GenerateNetworkBaseTestSuite) SetupTest() {
	suite.EnvBaseTestSuite.SetupTest()

	httpmock.ActivateNonDefault(env.HTTPClientV4)
	httpmock.ActivateNonDefault(env.HTTPClientV6)
}

func (suite *GenerateNetworkBaseTestSuite) TearDownTest() {
	httpmock.DeactivateAndReset()
	suite.EnvBaseTestSuite.TearDownTest()
}

type GenerateNetworkTestSuite struct {
	GenerateNetworkBaseTestSuite
}

func (suite *GenerateNetworkTestSuite) SetupResponder(r httpmock.Responder) {
	httpmock.RegisterResponder(http.MethodGet, "https://ipinfo.io/json", r)
}

func (suite *GenerateNetworkTestSuite) TestEnvSet() {
	suite.Setenv(env.EnvNetworkIPv4, "xx")
	env.GenerateNetwork(suite.ctx, suite.values, suite.wg)
	suite.Empty(suite.Collect())
}

func (suite *GenerateNetworkTestSuite) TestCannotAccess() {
	suite.SetupResponder(httpmock.NewErrorResponder(io.ErrUnexpectedEOF))
	env.GenerateNetwork(suite.ctx, suite.values, suite.wg)
	suite.Empty(suite.Collect())
}

func (suite *GenerateNetworkTestSuite) TestBadHTTPStatusCode() {
	suite.SetupResponder(httpmock.NewBytesResponder(http.StatusBadGateway, nil))
	env.GenerateNetwork(suite.ctx, suite.values, suite.wg)
	suite.Empty(suite.Collect())
}

func (suite *GenerateNetworkTestSuite) TestBrokenJSON() {
	suite.SetupResponder(httpmock.NewBytesResponder(http.StatusOK, nil))
	env.GenerateNetwork(suite.ctx, suite.values, suite.wg)
	suite.Empty(suite.Collect())
}

func (suite *GenerateNetworkTestSuite) TestCorrectResponse() {
	suite.SetupResponder(httpmock.NewJsonResponderOrPanic(
		http.StatusOK,
		map[string]string{
			"ip":       "127.0.0.1",
			"hostname": "hostname.provider.com",
			"city":     "XXX",
			"region":   "RRR",
			"country":  "CC",
			"loc":      "1234.56,78.9",
			"org":      "AS0000 OOO",
			"postal":   "0000",
			"timezone": "Europe/Berlin",
			"readme":   "https://ipinfo.io/missingauth",
		}))
	env.GenerateNetwork(suite.ctx, suite.values, suite.wg)
	data := suite.Collect()

	suite.Len(data, 11)
	suite.Equal("127.0.0.1", data[env.EnvNetworkIPv4])
	suite.Equal("hostname.provider.com", data[env.EnvNetworkHostname])
	suite.Equal("XXX", data[env.EnvNetworkCity])
	suite.Equal("RRR", data[env.EnvNetworkRegion])
	suite.Equal("CC", data[env.EnvNetworkCountry])
	suite.Equal("1234.56", data[env.EnvNetworkLatitude])
	suite.Equal("78.9", data[env.EnvNetworkLongitude])
	suite.Equal("OOO", data[env.EnvNetworkOrganization])
	suite.Equal("0000", data[env.EnvNetworkASN])
	suite.Equal("0000", data[env.EnvNetworkPostal])
	suite.Equal("Europe/Berlin", data[env.EnvNetworkTimezone])
}

type GenerateNetworkIPv6TestSuite struct {
	GenerateNetworkBaseTestSuite
}

func (suite *GenerateNetworkIPv6TestSuite) SetupResponder(r httpmock.Responder) {
	httpmock.RegisterResponder(http.MethodGet, "https://ifconfig.co", r)
}

func (suite *GenerateNetworkIPv6TestSuite) TestEnvSet() {
	suite.Setenv(env.EnvNetworkIPv6, "xx")
	env.GenerateNetworkIPv6(suite.ctx, suite.values, suite.wg)
	suite.Empty(suite.Collect())
}

func (suite *GenerateNetworkIPv6TestSuite) TestCannotAccess() {
	suite.SetupResponder(httpmock.NewErrorResponder(io.ErrUnexpectedEOF))
	env.GenerateNetworkIPv6(suite.ctx, suite.values, suite.wg)
	suite.Empty(suite.Collect())
}

func (suite *GenerateNetworkIPv6TestSuite) TestBadHTTPStatusCode() {
	suite.SetupResponder(httpmock.NewBytesResponder(http.StatusBadGateway, nil))
	env.GenerateNetworkIPv6(suite.ctx, suite.values, suite.wg)
	suite.Empty(suite.Collect())
}

func (suite *GenerateNetworkIPv6TestSuite) TestBrokenJSON() {
	suite.SetupResponder(httpmock.NewBytesResponder(http.StatusOK, nil))
	env.GenerateNetworkIPv6(suite.ctx, suite.values, suite.wg)
	suite.Empty(suite.Collect())
}

func (suite *GenerateNetworkIPv6TestSuite) TestCorrectResponse() {
	suite.SetupResponder(httpmock.NewJsonResponderOrPanic(
		http.StatusOK,
		map[string]string{"ip": "cafe::1"}))
	env.GenerateNetworkIPv6(suite.ctx, suite.values, suite.wg)
	data := suite.Collect()

	suite.Len(data, 1)
	suite.Equal("cafe::1", data[env.EnvNetworkIPv6])
}

func TestGenerateNetwork(t *testing.T) {
	suite.Run(t, &GenerateNetworkTestSuite{})
}

func TestGenerateNetworkIPv6(t *testing.T) {
	suite.Run(t, &GenerateNetworkIPv6TestSuite{})
}
