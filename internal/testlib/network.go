package testlib

import (
	"testing"

	"github.com/9seconds/chore/internal/network"
	"github.com/jarcoal/httpmock"
)

type NetworkTestSuite struct {
	t      *testing.T
	dialer *DialerMock
	dns    *DNSResolverMock
}

func (suite *NetworkTestSuite) Setup(t *testing.T) {
	t.Helper()

	suite.t = t

	httpmock.Activate()
	httpmock.ActivateNonDefault(network.HTTPClient)
	httpmock.ActivateNonDefault(network.HTTPClientV4)
	httpmock.ActivateNonDefault(network.HTTPClientV6)
	t.Cleanup(httpmock.DeactivateAndReset)

	suite.dialer = &DialerMock{}
	oldDialer := network.NetDialer
	network.NetDialer = suite.dialer

	suite.dns = &DNSResolverMock{}
	oldDNS := network.DNSResolver
	network.DNSResolver = suite.dns

	t.Cleanup(func() {
		suite.dialer.AssertExpectations(t)
		suite.dns.AssertExpectations(t)
		network.NetDialer = oldDialer
		network.DNSResolver = oldDNS
	})
}

func (suite *NetworkTestSuite) Dialer() *DialerMock {
	return suite.dialer
}

func (suite *NetworkTestSuite) DNS() *DNSResolverMock {
	return suite.dns
}

func (suite *NetworkTestSuite) MakeNetConn() *NetConnMock {
	rv := &NetConnMock{}

	suite.t.Cleanup(func() {
		rv.AssertExpectations(suite.t)
	})

	return rv
}
