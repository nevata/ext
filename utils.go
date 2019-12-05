package ext

import (
	"os"
	"path/filepath"
)

var (
	ExeDir string
)

func init() {
	s, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	ExeDir = s
}

func FileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
