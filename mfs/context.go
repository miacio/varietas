package mfs

import (
	"sync"

	"github.com/miacio/varietas/util"
)

type Mfs interface {
	Next(...int)            // Next
	Back(...int)            // Back
	Stop()                  // Stop
	Start()                 // Start
	Close(error)            // Close return
	Error() (string, error) // Error
	Get(string) any         // Get
	Set(string, any)        // Set
}

var (
	_ Mfs = (*Context)(nil)
	_ Mfs = (*Factory)(nil)
)

// Context
type Context struct {
	*Factory
	params   map[string]any // params
	wg       sync.WaitGroup // method loading group
	stop     sync.WaitGroup // is stop
	next     chan int       // next index
	TaskName string         // taskName
	Now      int            // now index
	err      error          // context error
}

func (ctx *Context) Next(index ...int) {
	next := 1
	if index != nil {
		next = util.PositiveInt(index[0])
	}
	ctx.next <- ctx.Now + next
}

func (ctx *Context) Back(index ...int) {
	back := 1
	if index != nil {
		back = util.PositiveInt(index[0])
	}
	ctx.next <- ctx.Now - back
}

func (ctx *Context) Stop() {
	ctx.stop.Add(1)
	ctx.stop.Wait()
}

func (ctx *Context) Start() {
	ctx.stop.Done()
}

func (ctx *Context) Close(err error) {
	ctx.next <- -1
	ctx.err = err
}

func (ctx *Context) Error() (string, error) {
	return ctx.TaskName, ctx.err
}

func (ctx *Context) Get(name string) any {
	return ctx.params[name]
}

func (ctx *Context) Set(name string, val any) {
	ctx.params[name] = val
}
