package update_test

import (
	"errors"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/9seconds/chore/internal/testlib"
	"github.com/9seconds/chore/internal/update"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MetadataTestSuite struct {
	UpdateTestSuite
}

func (suite *MetadataTestSuite) TestBadPaging() {
	pages := []int{1, 2}
	responders := map[string]httpmock.Responder{
		"err":        httpmock.NewErrorResponder(errors.New("err")),
		"bad_status": httpmock.NewStringResponder(http.StatusBadRequest, "[]"),
		"bad_json":   httpmock.NewStringResponder(http.StatusOK, ""),
		"bad_reader": testlib.NetworkResponderFromReader(
			http.StatusOK,
			testlib.NetworkBrokenReader()),
	}
	withUnstable := []bool{true, false}

	for _, pageNumber := range pages {
		for responderName, responder := range responders {
			for _, unstable := range withUnstable {
				pageNumber := pageNumber
				responderName := responderName
				responder := responder
				unstable := unstable

				suite.T().Run(strconv.Itoa(pageNumber), func(t *testing.T) {
					t.Run(responderName, func(t *testing.T) {
						t.Run(strconv.FormatBool(unstable), func(t *testing.T) {
							if pageNumber == 1 {
								suite.Register1PageResponder(responder)
							} else {
								suite.Register2PageResponder(responder)
							}

							_, err := update.GetLatestRelease(suite.Context(), unstable)
							assert.Error(t, err)
						})
					})
				})
			}
		}
	}
}

func (suite *MetadataTestSuite) TestBadFiles() {
	files := []string{"metadata", "artifacts"}
	responders := map[string]httpmock.Responder{
		"_err":        httpmock.NewErrorResponder(errors.New("err")),
		"_bad_status": httpmock.NewStringResponder(http.StatusBadRequest, "[]"),
		"_bad_json":   httpmock.NewStringResponder(http.StatusOK, ""),
		"_bad_reader": testlib.NetworkResponderFromReader(
			http.StatusOK,
			testlib.NetworkBrokenReader()),
		"metadata_bad_project": httpmock.NewJsonResponderOrPanic(
			http.StatusOK,
			suite.GetFixture("metadata_bad_project.json")),
		"artifacts_bad_1": httpmock.NewJsonResponderOrPanic(
			http.StatusOK,
			[]string{}),
		"artifacts_bad_2": httpmock.NewJsonResponderOrPanic(
			http.StatusOK,
			[]map[string]any{
				{
					"name":          "n1",
					"goos":          runtime.GOOS,
					"goarch":        "",
					"internal_type": 1,
				},
				{
					"name":          "n2",
					"goos":          "",
					"goarch":        "",
					"internal_type": 1,
				},
				{
					"name":          "n3",
					"goos":          "",
					"goarch":        runtime.GOARCH,
					"internal_type": 1,
				},
				{
					"name":          "n4",
					"goos":          runtime.GOOS,
					"goarch":        runtime.GOARCH,
					"internal_type": 15,
				},
			}),
		"artifacts_bad_3": httpmock.NewJsonResponderOrPanic(
			http.StatusOK,
			[]map[string]any{
				{
					"name":          "n1",
					"goos":          runtime.GOOS,
					"goarch":        runtime.GOARCH,
					"internal_type": 1,
				},
				{
					"name":          "n2",
					"goos":          runtime.GOOS,
					"goarch":        runtime.GOARCH,
					"internal_type": 1,
				},
				{
					"name":          "n2",
					"goos":          runtime.GOOS,
					"goarch":        runtime.GOARCH,
					"internal_type": 13,
				},
			}),
		"artifacts_bad_4": httpmock.NewJsonResponderOrPanic(
			http.StatusOK,
			[]map[string]any{
				{
					"name":          "n1",
					"goos":          runtime.GOOS,
					"goarch":        runtime.GOARCH,
					"internal_type": 1,
				},
				{
					"name":          "n2",
					"goos":          runtime.GOOS,
					"goarch":        runtime.GOARCH + "xxx",
					"internal_type": 1,
				},
			}),
		"artifacts_bad_5": httpmock.NewJsonResponderOrPanic(
			http.StatusOK,
			[]map[string]any{
				{
					"name":          "n1",
					"goos":          runtime.GOOS,
					"goarch":        runtime.GOARCH,
					"internal_type": 1,
				},
				{
					"name":          "n2",
					"goos":          runtime.GOOS,
					"goarch":        runtime.GOARCH + "xxx",
					"internal_type": 1,
				},
				{
					"name":          "n",
					"goos":          runtime.GOOS,
					"goarch":        runtime.GOARCH,
					"internal_type": 13,
				},
			}),
	}
	withUnstable := []bool{true, false}

	for _, filename := range files {
		for _, unstable := range withUnstable {
			for responderName, responder := range responders {
				filename := filename
				unstable := unstable
				responderName := responderName
				responder := responder

				if !(strings.HasPrefix(responderName, "_") || strings.HasPrefix(responderName, filename)) {
					continue
				}

				suite.T().Run(filename, func(t *testing.T) {
					t.Run(strconv.FormatBool(unstable), func(t *testing.T) {
						t.Run(responderName, func(t *testing.T) {
							if filename == "metadata" {
								suite.RegisterMetadataResponder(responder)
							} else {
								suite.RegisterArtifactsResponder(responder)
							}

							_, err := update.GetLatestRelease(suite.Context(), unstable)
							assert.Error(t, err)
						})
					})
				})
			}
		}
	}
}

func (suite *MetadataTestSuite) TestOk() {
	release, err := update.GetLatestRelease(suite.Context(), true)
	suite.NoError(err)
	suite.Equal("0.0.1", release.Version)
	suite.Equal(URLArchive, release.ArchiveURL)
	suite.Equal(URLSignature, release.SignatureURL)
}

func TestMetadata(t *testing.T) {
	suite.Run(t, &MetadataTestSuite{})
}
