package v1_test

import (
	"bytes"
	"sort"
	"testing"

	"github.com/9seconds/chore/internal/testlib"
	v1 "github.com/9seconds/chore/internal/vault/v1"
	"github.com/stretchr/testify/suite"
)

type VaultTestSuite struct {
	suite.Suite

	testlib.FixturesTestSuite
}

func (suite *VaultTestSuite) SetupTest() {
	suite.FixturesTestSuite.Setup(suite.T())
}

func (suite *VaultTestSuite) TestNewWithEmptyPassword() {
	_, err := v1.NewVault("")
	suite.ErrorIs(err, v1.ErrEmptyPassword)
}

func (suite *VaultTestSuite) TestNoKey() {
	vault, err := v1.NewVault("pass")
	suite.NoError(err)

	_, ok := vault.Get("k1")
	suite.False(ok)
}

func (suite *VaultTestSuite) TestVersion() {
	vault, err := v1.NewVault("pass")
	suite.NoError(err)

	suite.EqualValues(1, vault.Version())
}

func (suite *VaultTestSuite) TestSetGetDel() {
	vault, err := v1.NewVault("pass")
	suite.NoError(err)

	vault.Set("k1", "v1")

	value, exists := vault.Get("k1")
	suite.True(exists)
	suite.Equal("v1", value)

	vault.Delete("k1")

	_, exists = vault.Get("k1")
	suite.False(exists)
}

func (suite *VaultTestSuite) TestDoubleGet() {
	vault, err := v1.NewVault("pass")
	suite.NoError(err)

	vault.Set("k2", "v2")

	vault.Delete("k2")
	vault.Delete("k2")
	vault.Delete("k1")
}

func (suite *VaultTestSuite) TestList() {
	vault, err := v1.NewVault("pass")
	suite.NoError(err)

	suite.Empty(vault.List())

	vault.Set("k1", "v1")
	suite.Equal([]string{"k1"}, vault.List())

	vault.Set("k2", "v1")

	list := vault.List()
	sort.Strings(list)
	suite.Equal([]string{"k1", "k2"}, list)
}

func (suite *VaultTestSuite) TestRestore() {
	vault, err := v1.NewVault("pass")
	suite.NoError(err)

	vault.Set("k1", "v1")
	vault.Set("k2", "v2")

	data, err := vault.MarshalBinary()
	suite.NoError(err)

	suite.NotContains([]byte("pass"), data)
	suite.NotContains([]byte("k1"), data)
	suite.NotContains([]byte("k2"), data)

	newVault, err := v1.NewVault("pass")
	suite.NoError(err)

	suite.NoError(newVault.UnmarshalBinary(data))

	value, ok := newVault.Get("k1")
	suite.True(ok)
	suite.Equal("v1", value)

	value, ok = newVault.Get("k2")
	suite.True(ok)
	suite.Equal("v2", value)
}

func (suite *VaultTestSuite) TestBadRestore() {
	vault, err := v1.NewVault("pass")
	suite.NoError(err)

	vault.Set("k1", "v1")
	vault.Set("k2", "v2")

	data, err := vault.MarshalBinary()
	suite.NoError(err)

	newVault, err := v1.NewVault("pass2")
	suite.NoError(err)

	suite.ErrorIs(newVault.UnmarshalBinary(data), v1.ErrBadPassword)
}

func (suite *VaultTestSuite) TestCannotReadKDFNonce() {
	vault, err := v1.NewVault("pass")
	suite.NoError(err)

	suite.ErrorIs(
		vault.UnmarshalBinary(bytes.Repeat([]byte{0}, v1.NonceLength-1)),
		v1.ErrShortData)
}

func (suite *VaultTestSuite) TestCannotReadMAC() {
	vault, err := v1.NewVault("pass")
	suite.NoError(err)

	suite.ErrorIs(
		vault.UnmarshalBinary(
			bytes.Repeat([]byte{0}, v1.NonceLength+1)),
		v1.ErrShortData)
}

func (suite *VaultTestSuite) TestCannotReadLength() {
	vault, err := v1.NewVault("pass")
	suite.NoError(err)

	suite.ErrorIs(
		vault.UnmarshalBinary(
			bytes.Repeat([]byte{0}, v1.NonceLength+v1.MACLength+1)),
		v1.ErrShortData)
}

func (suite *VaultTestSuite) TestCannotReadMessageLength() {
	vault, err := v1.NewVault("pass")
	suite.NoError(err)

	suite.ErrorIs(
		vault.UnmarshalBinary(
			bytes.Repeat(
				[]byte{1}, v1.NonceLength+v1.MACLength+v1.LenLength)),
		v1.ErrShortData)
}

func (suite *VaultTestSuite) TestReadCorrectSnapshot() {
	vault, err := v1.NewVault("pass")
	suite.NoError(err)

	vault.Set("k1", "v1")

	marshalled, _ := vault.MarshalBinary()
	suite.EnsureSnapshot(marshalled, "correct-snapshot")

	data := suite.ReadPath("correct-snapshot")
	suite.NoError(vault.UnmarshalBinary(data))

	value, ok := vault.Get("k1")
	suite.True(ok)
	suite.Equal("v1", value)
}

func (suite *VaultTestSuite) TestSnapshotWithWrongPassword() {
	vault, err := v1.NewVault("pass")
	suite.NoError(err)

	vault.Set("k1", "v1")

	marshalled, _ := vault.MarshalBinary()
	suite.EnsureSnapshot(marshalled, "incorrect-snapshot")

	vault, err = v1.NewVault("pass2")
	suite.NoError(err)

	data := suite.ReadPath("incorrect-snapshot")
	suite.ErrorIs(vault.UnmarshalBinary(data), v1.ErrBadPassword)
}

func TestVault(t *testing.T) {
	suite.Run(t, &VaultTestSuite{})
}
