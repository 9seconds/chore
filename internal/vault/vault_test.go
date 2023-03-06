package vault_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"testing/iotest"

	"github.com/9seconds/chore/internal/testlib"
	"github.com/9seconds/chore/internal/vault"
	v1 "github.com/9seconds/chore/internal/vault/v1"
	"github.com/stretchr/testify/suite"
)

type BaseVaultTest struct {
	suite.Suite
}

func (suite *BaseVaultTest) TestLatestVersion() {
	box, err := vault.New("pass")
	suite.NoError(err)
	suite.EqualValues(vault.LatestVersion, box.Version())
}

func (suite *BaseVaultTest) TestEmptyPassword() {
	_, err := vault.New("")
	suite.ErrorContains(err, "password")
}

func (suite *BaseVaultTest) TestSaveLoad() {
	box, err := vault.New("pass")
	suite.NoError(err)

	box.Set("k1", "v1")
	box.Set("k2", "v2")

	path := filepath.Join(suite.T().TempDir(), "file")
	fileHandler, err := os.Create(path)
	suite.NoError(err)

	suite.NoError(vault.Save(fileHandler, box))
	suite.NoError(fileHandler.Close())

	fileHandler, err = os.Open(path)
	suite.NoError(err)

	box2, err := vault.Open(fileHandler, "pass")
	suite.NoError(err)

	value, exists := box2.Get("k1")
	suite.True(exists)
	suite.Equal("v1", value)

	value, exists = box2.Get("k2")
	suite.True(exists)
	suite.Equal("v2", value)
}

func (suite *BaseVaultTest) TestIncorrectPassword() {
	box, err := vault.New("pass")
	suite.NoError(err)

	path := filepath.Join(suite.T().TempDir(), "file")
	fileHandler, err := os.Create(path)
	suite.NoError(err)

	suite.NoError(vault.Save(fileHandler, box))
	suite.NoError(fileHandler.Close())

	fileHandler, err = os.Open(path)
	suite.NoError(err)

	_, err = vault.Open(fileHandler, "bad-password")
	suite.ErrorContains(err, "password")
}

func (suite *BaseVaultTest) TestEmptyData() {
	_, err := vault.Open(&bytes.Buffer{}, "pass")
	suite.ErrorContains(err, "cannot read version")
}

func (suite *BaseVaultTest) TestCannotReadData() {
	data := bytes.NewBuffer([]byte{1})

	_, err := vault.Open(iotest.TimeoutReader(data), "pass")
	suite.ErrorContains(err, "cannot read data")
}

func (suite *BaseVaultTest) TestUnsupportedVersion() {
	data := bytes.NewBuffer([]byte{0, 0})

	_, err := vault.Open(data, "pass")
	suite.ErrorIs(err, vault.ErrUnsupportedVaultVersion)
}

func (suite *BaseVaultTest) TestCannotUnmarshalVault() {
	data := bytes.NewBuffer([]byte{byte(vault.LatestVersion)})

	_, err := vault.Open(data, "pass")
	suite.ErrorContains(err, "cannot unmarshal vault")
}

type VaultV1Test struct {
	suite.Suite

	testlib.FixturesTestSuite
}

func (suite *VaultV1Test) SetupTest() {
	suite.FixturesTestSuite.Setup(suite.T())

	box, err := v1.NewVault("pass")
	suite.NoError(err)

	box.Set("k1", "vv1")

	data := &bytes.Buffer{}
	suite.NoError(vault.Save(data, box))

	suite.EnsureSnapshot(data.Bytes(), "v1-snapshot")
}

func (suite *VaultV1Test) TestReadOk() {
	fp, err := os.Open(suite.FixturePath("v1-snapshot"))
	suite.NoError(err)

	defer fp.Close()

	box, err := vault.Open(fp, "pass")
	suite.NoError(err)

	value, ok := box.Get("k1")
	suite.True(ok)
	suite.Equal("vv1", value)
}

func (suite *VaultV1Test) TestReadFail() {
	fp, err := os.Open(suite.FixturePath("v1-snapshot"))
	suite.NoError(err)

	defer fp.Close()

	_, err = vault.Open(fp, "bad-password")
	suite.ErrorContains(err, "password")
}

func TestBaseVault(t *testing.T) {
	suite.Run(t, &BaseVaultTest{})
}

func TestVaultV1(t *testing.T) {
	suite.Run(t, &VaultV1Test{})
}
