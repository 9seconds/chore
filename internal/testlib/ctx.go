package testlib

import (
	"context"
	"testing"
)

type CtxTestSuite struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
}

func (suite *CtxTestSuite) Setup(t *testing.T) {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())

	suite.ctx = ctx
	suite.ctxCancel = cancel

	t.Cleanup(cancel)
}

func (suite *CtxTestSuite) Context() context.Context {
	return suite.ctx
}

func (suite *CtxTestSuite) ContextCancel() context.CancelFunc {
	return suite.ctxCancel
}
