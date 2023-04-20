package mfs

import "errors"

type TaskMethod struct {
	TaskName   string // taskName
	TaskMethod Method // TaskMethod
}

func (t *TaskMethod) Name() string {
	return t.TaskName
}

func (t *TaskMethod) Runner(ctx *Context) {
	t.TaskMethod(ctx)
}

func CreateMethod(name string, taskMethod func(ctx *Context)) TaskMethod {
	return TaskMethod{
		TaskName:   name,
		TaskMethod: taskMethod,
	}
}

type Method func(*Context)

type MethodChain []TaskMethod

type Factory struct {
	*Context
	methodChain MethodChain // MethodChain
}

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

func (f *Factory) AddTaskMethod(taskMethods ...TaskMethod) {
	f.methodChain = append(f.methodChain, taskMethods...)
}

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

func (f *Factory) Run() {
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
					f.TaskName = f.methodChain[ext].TaskName
					go f.methodChain[ext].TaskMethod(f.Context)
				}
			}
		}
	}()
	f.Now = 0
	f.next <- 0
	f.wg.Wait()
}
