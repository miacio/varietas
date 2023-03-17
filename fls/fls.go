package fls

import "os"

type context struct {
	file *os.File
}

// Open
func Open(name string) (*context, error) {
	return OpenFile(name, os.O_RDONLY, 0)
}

// OpenFile
func OpenFile(name string, flag int, perm os.FileMode) (*context, error) {
	file, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	return &context{file: file}, nil
}

// Create
func Create(name string) (*context, error) {
	return OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
}
