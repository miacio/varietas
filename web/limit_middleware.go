package web

import "github.com/gin-gonic/gin"

// Limit
// gin-limit stay under the limit with this handy gin middleware
func Limit(n int) gin.HandlerFunc {
	sem := make(chan struct{}, n)
	acquire := func() { sem <- struct{}{} }
	release := func() { <-sem }
	return func(c *gin.Context) {
		acquire()       // before request
		defer release() // after request
		c.Next()
	}
}
