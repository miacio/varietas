package web_test

import (
	"net/http"
	"testing"

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
