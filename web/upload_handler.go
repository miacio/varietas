package web

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type ChunkFileRequest struct {
	FileId   string   `json:"fileId"`   // client create uuid
	FileName string   `json:"fileName"` // file name
	FileKeys []string `json:"fileKeys"` // file slice all key md5
	FileKey  string   `json:"fileKey"`  // file now key to md5 - if server read the slice to md5 eq key not eq then fail

	File *multipart.FileHeader `json:"file"` // now file

	ctx *gin.Context // ctx
}

func (cf *ChunkFileRequest) BindingForm(c *gin.Context) error {
	cf.FileId = c.PostForm("fileId")
	cf.FileName = c.PostForm("fileName")
	cf.FileKeys = c.PostFormArray("fileKeys")
	cf.FileKey = c.PostForm("fileKey")
	cf.ctx = c
	upFile, err := c.FormFile("file")

	if err != nil {
		return err
	}
	cf.File = upFile
	return cf.md5()
}

func (cf *ChunkFileRequest) md5() error {
	muFile, err := cf.File.Open()
	if err != nil {
		return err
	}
	defer muFile.Close()
	bt, err := io.ReadAll(muFile)
	if err != nil {
		return err
	}
	hash := fmt.Sprintf("%x", md5.Sum(bt))
	if hash != cf.FileKey {
		return errors.New("current file slice key error")
	}
	return nil
}

func (cf *ChunkFileRequest) SaveUploadedFile(tempPath, path string) error {
	tempFolder := filepath.Join(tempPath, cf.FileId)

	_, err := os.Stat(tempFolder)
	if os.IsNotExist(err) {
		err := os.MkdirAll(tempFolder, os.ModePerm)
		if err != nil {
			return err
		}
	}

	if err := cf.ctx.SaveUploadedFile(cf.File, filepath.Join(tempFolder, cf.FileKey)); err != nil {
		return err
	}

	for _, fileKey := range cf.FileKeys {
		tempFile := filepath.Join(tempFolder, fileKey)
		if _, err := os.Stat(tempFile); err != nil {
			return nil
		}
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0664)
	if err != nil {
		return err
	}

	defer file.Close()

	for _, fileKey := range cf.FileKeys {
		tempFile := filepath.Join(tempFolder, fileKey)
		bt, err := os.ReadFile(tempFile)
		if err != nil {
			return err
		}
		file.Write(bt)
	}

	return os.RemoveAll(tempFolder)
}

// param: fileId
// param: fileName
// param: fileKeys the file slice all file key md5
// param: fileKey  now file slice key md5
// param: file     now slice file
func ChunkFile(c *gin.Context) {
	var cf ChunkFileRequest

	if err := cf.BindingForm(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "400", "msg": "bad file param", "err": err.Error()})
		return
	}

	if err := cf.SaveUploadedFile("./temp", "./uploads/"+cf.FileName); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"code": "503", "msg": "bad save upload file", "err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": "200", "msg": "success"})
}
