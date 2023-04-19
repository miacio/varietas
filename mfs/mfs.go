package mfs

import (
	"errors"
	"sync"
)

// Context
type Context struct {
	TaskId string         // now task id
	next   chan int       // ctx.Next method to chan
	stop   bool           // ctx is stop
	now    int            // now index
	params map[string]any // context core message
	err    error          // result err
}

func (c *Context) Back(a ...int) {
	if a == nil {
		c.next <- c.now - 1
	} else {
		if a[0] < 0 {
			a[0] = a[0] - a[0] - a[0]
		}
		c.next <- c.now - a[0]
	}
}

func (c *Context) Next() {
	c.next <- c.now + 1
}

func (c *Context) Close(err error) error {
	c.next <- -1
	c.err = err
	return err
}

func (c *Context) Stop() {
	c.stop = true
}

func (c *Context) Set(key string, val any) {
	c.params[key] = val
}

func (c *Context) Get(key string) any {
	return c.params[key]
}

// Method
type Method interface {
	TaskId() string  // return this method name
	Do(ctx *Context) // runner func
}

// MethodChain
type MethodChain []Method

func (m MethodChain) Length() int {
	return len(m)
}

type Factory struct {
	ctx           *Context    // context
	methods       MethodChain // factory method chain is method slice
	appendMethods MethodChain // methods added when running the runner
	runner        bool        // runner
	wg            sync.WaitGroup
}

// GenerateFactory
func GenerateFactory(name string) *Factory {
	f := &Factory{
		ctx: &Context{
			next:   make(chan int, 1),
			now:    0,
			params: make(map[string]any, 0),
			err:    nil,
		},
		methods: make([]Method, 0),
	}
	f.ctx.next <- 0
	return f
}

func (f *Factory) AppendMethods(methods ...Method) {
	if f.runner {
		f.appendMethods = append(f.appendMethods, methods...)
		return
	}
	f.methods = append(f.methods, methods...)
}

func (f *Factory) appendRunner() {
	if f.runner {
		return
	}
	f.methods = append(f.methods, f.appendMethods...)
}

// Factory.Do runner
func (f *Factory) do() {
	if f.methods == nil || f.methods.Length() == 0 {
		f.ctx.err = errors.New("factory in methods is empty")
		return
	}
	f.runner = true
	next := <-f.ctx.next
	f.ctx.now = next
	if next >= f.methods.Length() || next < 0 {
		f.runner = false
		f.appendRunner()
		f.wg.Done()
		return
	}
	f.ctx.TaskId = f.methods[next].TaskId()
	f.methods[next].Do(f.ctx)
	f.do()
}

func (f *Factory) Do() {
	f.wg.Add(1)
	go func() {
		f.do()
	}()
}

// EndMessage
func (f *Factory) EndMessage() (string, error) {
	f.wg.Wait()
	return f.ctx.TaskId, f.ctx.err
}

// Next
func (f *Factory) Next() bool {
	if f.ctx.stop {
		f.ctx.next <- f.ctx.now + 1
		f.ctx.stop = false
		return true
	}
	return false
}
