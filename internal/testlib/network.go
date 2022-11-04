package testlib

import (
	"testing"

	"github.com/9seconds/chore/internal/env"
	"github.com/jarcoal/httpmock"
)

type NetworkTestSuite struct{}

func (suite *NetworkTestSuite) Setup(t *testing.T) {
	t.Helper()

	httpmock.Activate()
	httpmock.ActivateNonDefault(env.HTTPClientV4)
	httpmock.ActivateNonDefault(env.HTTPClientV6)

	t.Cleanup(httpmock.DeactivateAndReset)
}
