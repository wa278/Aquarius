
package lib

import (
	"os"
	"path/filepath"
)

type PathConst struct {
	RootDir string
	DS      string

	ConfigDir     string
}
//获取路径信息
func NewPathConst() *PathConst {
	filename, _ := filepath.Abs(os.Args[0])
	rootDir := filepath.Dir(filepath.Dir(filename))

	ds := string(filepath.Separator)

	return &PathConst{
		RootDir:       rootDir,
		DS:            ds,
		ConfigDir:     rootDir + ds + "configs",
	}
}

