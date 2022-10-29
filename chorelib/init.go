package chorelib

import (
	"path/filepath"

	"github.com/adrg/xdg"
)

var (
	Home           = filepath.Join(xdg.ConfigHome, "chore")
	PersistentDirs = filepath.Join(xdg.DataHome, "chore")
)
