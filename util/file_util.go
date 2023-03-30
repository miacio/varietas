package util

import (
	"os"
	"path/filepath"
)

// FileExist check file exist returnes true, else false.
func FileExist(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// FileIsDir check file exist and file is folder
func FileIsDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}

// FolderIsNotExistThenMkdir
// check path folder is exist, if not exist then mkdirall the path
// else os.Stat(path) os.IsNotExist false, then return err, else return nil
// so used the method then check this err == nil then success.
func FolderIsNotExistThenMkdir(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(path, os.ModePerm)
		}
		return err
	}
	return nil
}

// FileFindAllFileChildren
// suffix is "" then
// search the path children all files
// suffix is not "" then
// search the path children files suffix is param suffix
func FileFindAllFileChildren(path, suffix string) ([]string, error) {
	isDir, err := FileIsDir(path)
	if err != nil {
		return nil, err
	}
	if !isDir {
		fileNameWithSuffix := filepath.Base(path)
		fileSuffix := filepath.Ext(fileNameWithSuffix)
		if suffix == "" || fileSuffix == suffix {
			return []string{path}, nil
		}
		return nil, nil
	}
	dirEntrys, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0)
	for _, dirEntry := range dirEntrys {
		if dirEntry.IsDir() {
			files, err := FileFindAllFileChildren(path+"/"+dirEntry.Name(), suffix)
			if err != nil {
				continue
			}
			result = append(result, files...)
		} else {
			fileName := path + "/" + dirEntry.Name()
			fileNameWithSuffix := filepath.Base(fileName)
			fileSuffix := filepath.Ext(fileNameWithSuffix)
			if suffix == "" || fileSuffix == suffix {
				result = append(result, fileName)
			}
		}
	}
	return result, nil
}
