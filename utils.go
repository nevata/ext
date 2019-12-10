package ext

import (
	"os"
	"path/filepath"
)

//ExeDir 执行程序所在路径
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

//FileExist 文件和目录判断
func FileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
