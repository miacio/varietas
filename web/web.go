package web

import (
	"github.com/gin-gonic/gin"
)

type context struct {
	*gin.Engine
	routers Routers
}

type Router interface {
	Execute(c *gin.Engine)
}

type Routers []Router

// New
func New(eng *gin.Engine) *context {
	return &context{
		Engine:  eng,
		routers: make(Routers, 0),
	}
}

// RegisterController struct
func (c *context) Register(routers ...Router) {
	c.routers = append(c.routers, routers...)
}

// Prepare execute this method to register a route before starting running
func (c *context) Prepare() {
	for _, router := range c.routers {
		router.Execute(c.Engine)
	}
}
