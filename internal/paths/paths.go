package paths

import (
	"path/filepath"

	"github.com/adrg/xdg"
)

const (
	ChoreDir = "chore"
)

func ConfigRoot() string {
	return filepath.Join(xdg.ConfigHome, ChoreDir)
}

func ConfigNamespace(ns string) string {
	return filepath.Join(ConfigRoot(), ns)
}

func ConfigNamespaceScript(ns, script string) string {
	return filepath.Join(ConfigNamespace(ns), script)
}

func ConfigNamespaceScriptConfig(ns, script string) string {
	return ConfigNamespaceScript(ns, script) + ".toml"
}

func DataRoot() string {
	return filepath.Join(xdg.DataHome, ChoreDir)
}

func DataNamespace(ns string) string {
	return filepath.Join(DataRoot(), ns)
}

func DataNamespaceScript(ns, script string) string {
	return filepath.Join(DataNamespace(ns), script)
}

func CacheRoot() string {
	return filepath.Join(xdg.CacheHome, ChoreDir)
}

func CacheNamespace(ns string) string {
	return filepath.Join(CacheRoot(), ns)
}

func CacheNamespaceScript(ns, script string) string {
	return filepath.Join(CacheNamespace(ns), script)
}

func StateRoot() string {
	return filepath.Join(xdg.StateHome, ChoreDir)
}

func StateNamespace(ns string) string {
	return filepath.Join(StateRoot(), ns)
}

func StateNamespaceScript(ns, script string) string {
	return filepath.Join(StateNamespace(ns), script)
}