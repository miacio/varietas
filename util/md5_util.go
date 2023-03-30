package util

import (
	"crypto/md5"
	"fmt"
	"io"
)

// MD5
func MD5(bt []byte) string {
	return fmt.Sprintf("%x", md5.Sum(bt))
}

// MD5File
func MD5File(r io.Reader) string {
	bf := make([]byte, 4096)
	hashMd5 := md5.New()
	for {
		n, err := r.Read(bf)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil && err != io.EOF {
			break
		}
		hashMd5.Write(bf[:n])
	}
	return fmt.Sprintf("%x", hashMd5.Sum(nil))
}
