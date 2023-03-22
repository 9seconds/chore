package paths

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

const (
	CacheDirTagName    = "CACHEDIR.TAG"
	CacheDirTagContent = `Signature: 8a477f597d28d172789f06886806bc55
# This file is a cache directory tag created by (application name).
# For information about cache directory tags, see:
#	https://bford.info/cachedir/
`

	DirectoryPermission fs.FileMode = 0o755
	FilePermission      fs.FileMode = 0o644
)

func EnsureDir(path string) error {
	return os.MkdirAll(path, DirectoryPermission)
}

func EnsureFile(path, content string) error {
	if err := EnsureDir(filepath.Dir(path)); err != nil {
		return fmt.Errorf("cannot ensure parent directory: %w", err)
	}

	if err := os.WriteFile(path, []byte(content), FilePermission); err != nil {
		return fmt.Errorf("cannot write to file: %w", err)
	}

	return nil
}

func EnsureRoots(namespace, script string) error {
	if err := EnsureDir(ConfigNamespace(namespace)); err != nil {
		return fmt.Errorf("cannot create config root: %w", err)
	}

	if err := EnsureDir(DataNamespaceScript(namespace, script)); err != nil {
		return fmt.Errorf("cannot create data root: %w", err)
	}

	if err := EnsureDir(CacheNamespaceScript(namespace, script)); err != nil {
		return fmt.Errorf("cannot create cache root: %w", err)
	}

	if err := EnsureDir(StateNamespaceScript(namespace, script)); err != nil {
		return fmt.Errorf("cannot create state root: %w", err)
	}

	if err := EnsureFile(CacheDirTagPath(), CacheDirTagContent); err != nil {
		return fmt.Errorf("cannot create cachedirtag: %w", err)
	}

	return nil
}
