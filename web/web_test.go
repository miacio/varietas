package web_test

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/miacio/varietas/web"
)

type defaultCtr struct{}

var DefaultCtr web.Router = (*defaultCtr)(nil)

func (*defaultCtr) Execute(c *gin.Engine) {
	c.GET("/test", func(ctx *gin.Context) {
		ctx.JSONP(http.StatusOK, gin.H{"message": "success"})
	})
}

func TestWeb001(t *testing.T) {
	w := web.New(gin.Default())
	w.Register(DefaultCtr)
	w.Prepare()
	w.Run(":8080")
}

// cd web dir
// go test-v -run TestWeb002
// you need Ctrl+C close the method
func TestWeb002(t *testing.T) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	w := web.New(gin.Default())
	w.Register(DefaultCtr)
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
