package web_test

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/miacio/varietas/web"
)

type defaultCtr struct{}

var DefaultCtr web.Router = (*defaultCtr)(nil)

func (*defaultCtr) Execute(c *gin.Engine) {
	c.GET("/test", func(ctx *gin.Context) {
		ctx.JSONP(http.StatusOK, gin.H{"message": "success"})
	})
}

type fileUploadCtr struct{}

var FileUploadCtr web.Router = (*fileUploadCtr)(nil)

func (*fileUploadCtr) Execute(c *gin.Engine) {
	c.POST("/chunkFile", web.ChunkFile)
}

func TestWeb001(t *testing.T) {
	w := web.New(gin.Default())
	w.Register(DefaultCtr)
	w.Prepare()
	w.Run(":8080")
}

// cd web dir
// go test -v -run TestWeb002
// you need Ctrl+C close the method
func TestChunkFileUploadServer(t *testing.T) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	w := web.New(gin.Default())
	w.Register(FileUploadCtr)
	w.Prepare()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: w,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

func TestChunkFileUploadClient(t *testing.T) {
	// your client file path
	filePath := ""
	fileName := filepath.Base(filePath)

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Fatalf("file stat fail: %v\n", err)
		return
	}

	const chunkSize = 1 << (10 * 2) * 30

	num := math.Ceil(float64(fileInfo.Size()) / float64(chunkSize))

	fi, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("open file fail: %v\n", err)
		return
	}

	fileKeyMap := make(map[string][]byte, 0)
	fileKeys := make([]string, 0)

	for i := 1; i <= int(num); i++ {
		file := make([]byte, chunkSize)
		fi.Seek((int64(i)-1)*chunkSize, 0)
		if len(file) > int(fileInfo.Size()-(int64(i)-1)*chunkSize) {
			file = make([]byte, fileInfo.Size()-(int64(i)-1)*chunkSize)
		}
		fi.Read(file)

		key := fmt.Sprintf("%x", md5.Sum(file))
		fileKeyMap[key] = file
		fileKeys = append(fileKeys, key)
	}

	fileId := uuid.NewString()

	for _, key := range fileKeys {
		req := web.ChunkFileRequest{
			FileId:   fileId,
			FileName: fileName,
			FileKey:  key,
			FileKeys: fileKeys,
			File:     fileKeyMap[key],
		}
		body, _ := json.Marshal(req)

		res, err := http.Post("http://127.0.0.1:8080/chunkFile", "application/json", bytes.NewBuffer(body))

		if err != nil {
			log.Fatalf("http post fail: %v", err)
			return
		}
		defer res.Body.Close()
		msg, _ := io.ReadAll(res.Body)
		fmt.Println(string(msg))
	}
}
