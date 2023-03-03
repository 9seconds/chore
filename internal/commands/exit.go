package commands

import (
	"os"

	"github.com/9seconds/chore/internal/paths"
)

func Exit(code int) {
	paths.TempDirCleanup()
	os.Exit(code)
}
