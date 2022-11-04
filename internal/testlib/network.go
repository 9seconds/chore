package testlib

import (
	"testing"

	"github.com/9seconds/chore/internal/network"
	"github.com/jarcoal/httpmock"
)

type NetworkTestSuite struct{}

func (suite *NetworkTestSuite) Setup(t *testing.T) {
	t.Helper()

	httpmock.Activate()
	httpmock.ActivateNonDefault(network.HTTPClient)
	httpmock.ActivateNonDefault(network.HTTPClientV4)
	httpmock.ActivateNonDefault(network.HTTPClientV6)
	t.Cleanup(httpmock.DeactivateAndReset)
}
