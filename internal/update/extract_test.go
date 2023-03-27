package update_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/9seconds/chore/internal/testlib"
	"github.com/9seconds/chore/internal/update"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ExtractTestSuite struct {
	UpdateTestSuite
}

func (suite *ExtractTestSuite) Extract() (io.Reader, error) {
	return update.Extract(suite.Context(), URLArchive, URLSignature)
}

func (suite *ExtractTestSuite) TestBadResponse() {
	files := []string{"archive", "signature"}
	responders := map[string]httpmock.Responder{
		"err":        httpmock.NewErrorResponder(errors.New("err")),
		"bad_status": httpmock.NewStringResponder(http.StatusNotFound, ""),
		"bad_reader": testlib.NetworkResponderFromReader(
			http.StatusOK,
			testlib.NetworkBrokenReader()),
		"big_response": testlib.NetworkResponderFromReader(
			http.StatusOK,
			bytes.NewReader(bytes.Repeat([]byte{0}, update.MaxFileSize+1))),
	}

	for _, filename := range files {
		for responderName, responder := range responders {
			filename := filename
			responderName := responderName
			responder := responder

			suite.T().Run(filename, func(t *testing.T) {
				t.Run(responderName, func(t *testing.T) {
					if filename == "archive" {
						suite.RegisterArchiveResponder(responder)
					} else {
						suite.RegisterSignatureResponder(responder)
					}

					_, err := suite.Extract()

					assert.Error(t, err)
				})
			})
		}
	}
}

func (suite *ExtractTestSuite) TestBadSignature() {
	suite.RegisterSignatureResponder(
		httpmock.NewStringResponder(http.StatusOK, "1"))

	_, err := suite.Extract()
	suite.Error(err)
}

func (suite *ExtractTestSuite) TestBrokenArchive() {
	fixtures := []string{
		"archive_broken_gz.tar.gz",
		"archive_empty_valid_tar.tar.gz",
		"archive_missing_chore.tar.gz",
		"archive_bad_tar.tar.gz",
	}

	for _, fixture := range fixtures {
		fixture := fixture

		suite.T().Run(fixture, func(t *testing.T) {
			suite.RegisterArchiveResponder(
				httpmock.NewBytesResponder(
					http.StatusOK,
					suite.GetFixture(fixture).Bytes()))
			suite.RegisterSignatureResponder(
				httpmock.NewBytesResponder(
					http.StatusOK,
					suite.GetFixture(fixture+".sig").Bytes()))

			_, err := suite.Extract()
			suite.Error(err)
		})
	}
}

func (suite *ExtractTestSuite) TestOk() {
	reader, err := suite.Extract()
	suite.NoError(err)

	hasher := sha256.New()
	io.Copy(hasher, reader) //nolint: errcheck

	suite.Equal(
		"s/XEPVcARVJ4//YOH5hRRqKP5AwPeXRjAYy6GU4ViQY=",
		base64.StdEncoding.EncodeToString(hasher.Sum(nil)))
}

func TestExtract(t *testing.T) {
	suite.Run(t, &ExtractTestSuite{})
}
