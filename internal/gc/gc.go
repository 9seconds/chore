package gc

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/9seconds/chore/internal/paths"
	"github.com/9seconds/chore/internal/script"
	"github.com/tchap/go-patricia/v2/patricia"
)

func Collect(validScripts []*script.Script) ([]string, error) { //nolint: cyclop
	safeFiles := map[string]bool{
		paths.CacheDirTagPath(): true,
		paths.AppConfigPath():   true,
	}
	safePaths := patricia.NewTrie()

	for _, scr := range validScripts {
		log.Printf("add %q to safe files", scr.Path())

		safeFiles[scr.Path()] = true
		safeFiles[paths.ConfigNamespaceScriptVault(scr.Namespace)] = true

		safePaths.Set(patricia.Prefix(scr.Path()), true)
		safePaths.Set(patricia.Prefix(scr.DataPath()), true)
		safePaths.Set(patricia.Prefix(scr.CachePath()), true)
		safePaths.Set(patricia.Prefix(scr.StatePath()), true)

		if _, err := script.ValidateConfig(scr.ConfigPath()); err == nil {
			safeFiles[scr.ConfigPath()] = true
		} else {
			log.Printf("cannot add %q to safe files: %v", scr.ConfigPath(), err)
		}
	}

	collected := []string{}
	queue := NewListset()

	for _, v := range getRootPaths() {
		safePaths.Set(patricia.Prefix(filepath.Join(v, "\x00")), true)
		queue.Add(v)
	}

	queueIter := queue.Iter()

	for queueIter.Scan() {
		path := queueIter.Next()
		prefix := patricia.Prefix(path)

		switch {
		// exact match of the file or path we want to maintain
		case safeFiles[path], safePaths.Match(prefix):
			continue

		// this path is still a prefix of some safe path, so we need
		// to continue traversal
		case safePaths.MatchSubtree(prefix):
			files, err := os.ReadDir(path)

			switch {
			case errors.Is(err, fs.ErrNotExist):
				// it is possible that directory does not exist
				// like state path for a script that was never executed.
				continue
			case err != nil:
				return nil, fmt.Errorf("cannot read directory %s: %w", path, err)
			}

			for _, info := range files {
				queueIter.Add(filepath.Join(path, info.Name()))
			}

		// some directory that is not a prefix of a safe subdirectory
		default:
			collected = append(collected, path)
		}
	}

	sort.Strings(collected)

	return collected, nil
}

func Remove(paths []string) error {
	sort.Sort(sort.Reverse(sort.StringSlice(paths)))

	safePaths := patricia.NewTrie()

	for _, rootPath := range getRootPaths() {
		safePaths.Set(patricia.Prefix(rootPath), true)
	}

	queue := NewListset()

	for _, seed := range paths {
		queue.Add(seed)
	}

	queueIter := queue.Iter()

	for queueIter.Scan() {
		path := queueIter.Next()

		if safePaths.MatchSubtree(patricia.Prefix(path)) {
			continue
		}

		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("cannot remove %s: %w", path, err)
		}

		rootPath := filepath.Dir(path)
		content, err := os.ReadDir(rootPath)

		switch {
		case err != nil:
			return fmt.Errorf("cannot read directory %s: %w", path, err)
		case len(content) == 0:
			queueIter.Add(rootPath)
		}
	}

	return nil
}

func getRootPaths() []string {
	return []string{
		paths.ConfigRoot(),
		paths.DataRoot(),
		paths.CacheRoot(),
		paths.StateRoot(),
	}
}
