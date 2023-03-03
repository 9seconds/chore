package paths

import (
	"log"
	"os"
	"sync"
)

var (
	tempPaths []string
	tempMutex sync.Mutex
)

func TempDir() (string, error) {
	tempMutex.Lock()
	defer tempMutex.Unlock()

	path, err := os.MkdirTemp("", ChoreDir+"-")
	if err != nil {
		return "", err
	}

	tempPaths = append(tempPaths, path)

	return path, nil
}

func TempDirCleanup() {
	tempMutex.Lock()
	defer tempMutex.Unlock()

	for _, dir := range tempPaths {
		if err := os.RemoveAll(dir); err != nil {
			log.Printf("cannot remove temp path %s: %v", dir, err)
		}
	}

	tempPaths = tempPaths[:0]
}
