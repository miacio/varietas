package mfs

import "errors"

// ITaskMethod abstract
type ITaskMethod interface {
	Name() string    // get the task name
	Method(*Context) // task method
}

// Method
type Method func(*Context)

// MethodChain
type MethodChain []ITaskMethod

// Factory
type Factory struct {
	*Context
	methodChain MethodChain // MethodChain
}

// NewFactory
func NewFactory() *Factory {
	f := &Factory{
		methodChain: make(MethodChain, 0),
	}
	ctx := &Context{
		Factory: f,
		params:  make(map[string]any, 0),
		Now:     -1,
	}
	f.Context = ctx
	return f
}

// AddTaskMethod as a write method area,
// it is not recommended for developers to call this method when running the factory,
// otherwise uncontrollable results may occur
func (f *Factory) AddTaskMethod(taskMethods ...ITaskMethod) {
	f.methodChain = append(f.methodChain, taskMethods...)
}

// DropTaskByName as a deletion method area,
// it is not recommended for developers to call this method when running the factory,
// otherwise uncontrollable results may occur
func (f *Factory) DropTaskByName(name string) {
	newMethodChain := make(MethodChain, 0)
	for _, method := range f.methodChain {
		if method.Name() == name {
			continue
		}
		newMethodChain = append(newMethodChain, method)
	}
	f.methodChain = newMethodChain
}

// Excute
func (f *Factory) Excute() {
	f.next = make(chan int)
	go func() {
		defer close(f.next)
		f.wg.Add(1)
		for {
			select {
			case ext := <-f.next:
				if f.methodChain == nil {
					f.err = errors.New("method chain is empty")
					f.wg.Done()
					return
				}
				if ext < 0 || ext > len(f.methodChain) {
					f.wg.Done()
					return
				} else {
					f.Now = ext
					f.TaskName = f.methodChain[ext].Name()
					go f.methodChain[ext].Method(f.Context)
				}
			}
		}
	}()
	f.Now = 0
	f.next <- 0
	f.wg.Wait()
}
