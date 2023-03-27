package update_test

import (
	"net/http"
	"runtime"

	"github.com/9seconds/chore/internal/testlib"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"
)

const (
	URLPagination = "https://api.github.com/repos/9seconds/chore/releases"
	URLMetadata   = "https://github.com/9seconds/chore/releases/download/v0.0.1/metadata.json"
	URLArtifacts  = "https://github.com/9seconds/chore/releases/download/v0.0.1/artifacts.json"
	URLArchive    = "https://github.com/9seconds/chore/releases/download/v0.0.1/chore_0.0.1_linux_amd64.tar.gz"
	URLSignature  = "https://github.com/9seconds/chore/releases/download/v0.0.1/chore_0.0.1_linux_amd64.tar.gz.sig"
)

type UpdateTestSuite struct {
	suite.Suite

	testlib.FixturesTestSuite
	testlib.NetworkTestSuite
	testlib.CtxTestSuite
}

func (suite *UpdateTestSuite) SetupSuite() {
	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
		// I can't possibly pass metadata and artifacts for everything
		// let's concentrate on linux/amd64 only
		suite.T().Skip("Update test can run only on linux/amd64")
	}
}

func (suite *UpdateTestSuite) SetupTest() {
	suite.FixturesTestSuite.Setup(suite.T())
	suite.NetworkTestSuite.Setup(suite.T())
	suite.CtxTestSuite.Setup(suite.T())

	suite.Register1PageResponder(
		httpmock.NewJsonResponderOrPanic(
			http.StatusOK,
			httpmock.File(suite.FixturePath("releases_page_1.json"))))
	suite.Register2PageResponder(
		httpmock.NewJsonResponderOrPanic(
			http.StatusOK,
			httpmock.File(suite.FixturePath("releases_page_2.json"))))
	suite.RegisterMetadataResponder(
		httpmock.NewJsonResponderOrPanic(
			http.StatusOK,
			suite.GetFixture("metadata.json")))
	suite.RegisterArtifactsResponder(
		httpmock.NewJsonResponderOrPanic(
			http.StatusOK,
			suite.GetFixture("artifacts.json")))
	suite.RegisterArchiveResponder(
		httpmock.NewBytesResponder(
			http.StatusOK,
			suite.GetFixture("archive.tar.gz").Bytes()))
	suite.RegisterSignatureResponder(
		httpmock.NewBytesResponder(
			http.StatusOK,
			suite.GetFixture("archive.tar.gz.sig").Bytes()))
	httpmock.NewNotFoundResponder(suite.T().Fatal)
}

func (suite *UpdateTestSuite) GetFixture(name string) httpmock.File {
	return httpmock.File(suite.FixturePath(name))
}

func (suite *UpdateTestSuite) Register1PageResponder(resp httpmock.Responder) {
	httpmock.RegisterResponderWithQuery(
		http.MethodGet,
		URLPagination,
		"page=1&per_page=100",
		resp)
}

func (suite *UpdateTestSuite) Register2PageResponder(resp httpmock.Responder) {
	httpmock.RegisterResponderWithQuery(
		http.MethodGet,
		URLPagination,
		"page=2&per_page=100",
		resp)
}

func (suite *UpdateTestSuite) RegisterMetadataResponder(resp httpmock.Responder) {
	httpmock.RegisterResponder(http.MethodGet, URLMetadata, resp)
}

func (suite *UpdateTestSuite) RegisterArtifactsResponder(resp httpmock.Responder) {
	httpmock.RegisterResponder(http.MethodGet, URLArtifacts, resp)
}

func (suite *UpdateTestSuite) RegisterArchiveResponder(resp httpmock.Responder) {
	httpmock.RegisterResponder(http.MethodGet, URLArchive, resp)
}

func (suite *UpdateTestSuite) RegisterSignatureResponder(resp httpmock.Responder) {
	httpmock.RegisterResponder(http.MethodGet, URLSignature, resp)
}
