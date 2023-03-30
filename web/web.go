package web

import (
	"github.com/gin-gonic/gin"
	"github.com/miacio/varietas/util"
)

type Engine struct {
	*gin.Engine
	routers Routers
}

type Router interface {
	Execute(c *gin.Engine)
}

type Routers []Router

// New
func New(eng *gin.Engine) *Engine {
	return &Engine{
		Engine:  eng,
		routers: make(Routers, 0),
	}
}

// RegisterController struct
func (c *Engine) Register(routers ...Router) {
	c.routers = append(c.routers, routers...)
}

// Prepare execute this method to register a route before starting running
func (c *Engine) Prepare() {
	for _, router := range c.routers {
		router.Execute(c.Engine)
	}
}

// LoadHTMLFolders loads a slice of HTML folders in files
// and associate the result with HTML renderer.
// suffix is "" then load all folders in files other return this suffix files
func (c *Engine) LoadHTMLFolders(folders []string, suffix string) {
	if folders == nil {
		return
	}
	files := make([]string, 0)
	for _, folder := range folders {
		inFiles, err := util.FileFindAllFileChildren(folder, suffix)
		if err != nil {
			files = append(files, inFiles...)
		}
	}
	files = util.SliceDistinct(files...)
	if files != nil {
		c.LoadHTMLFiles(files...)
	}
}
