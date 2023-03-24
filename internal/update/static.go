package update

import (
	"embed"

	"github.com/ProtonMail/go-crypto/openpgp"
)

//go:embed static/*
var staticFS embed.FS

func getKeyring() openpgp.EntityList {
	file, err := staticFS.Open("static/keyring.asc")
	if err != nil {
		panic(err.Error())
	}

	keyring, err := openpgp.ReadArmoredKeyRing(file)
	if err != nil {
		panic(err.Error())
	}

	return keyring
}
