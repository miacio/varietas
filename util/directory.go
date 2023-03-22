package util

import (
	"os"
)

type directory struct {
	path string      //dir path
	info os.FileInfo // info
}

// Directory
// dir manager
func Directory(path string) *directory {
	dir := &directory{
		path: path,
	}
	info, err := os.Stat(path)
	if err == nil {
		dir.info = info
	}
	return dir
}

// IsExist
func (d *directory) IsExist() bool {
	return false
}

// IsDir
func (d *directory) IsDir() bool {
	if d.info == nil {
		return false
	}
	return d.info.IsDir()
}
