package testlib

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"testing"
	"testing/iotest"

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
	suite.t.Helper()

	return suite.dialer
}

func (suite *NetworkTestSuite) DNS() *DNSResolverMock {
	suite.t.Helper()

	return suite.dns
}

func (suite *NetworkTestSuite) MakeNetConn() *NetConnMock {
	suite.t.Helper()

	connMock := &NetConnMock{}

	suite.t.Cleanup(func() {
		connMock.AssertExpectations(suite.t)
	})

	return connMock
}

func NetworkResponderFromReader(status int, reader io.Reader) httpmock.Responder {
	resp := &http.Response{
		Status:        strconv.Itoa(status) + " " + http.StatusText(status),
		StatusCode:    status,
		Body:          io.NopCloser(reader),
		ContentLength: -1,
		Header:        http.Header{},
	}

	return httpmock.ResponderFromResponse(resp)
}

func NetworkBrokenReader() io.Reader {
	reader := iotest.OneByteReader(bytes.NewReader([]byte{1, 2, 3, 4}))
	reader = iotest.TimeoutReader(reader)

	return reader
}
