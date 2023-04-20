package mfs

// TaskMethod is ITaskMethod a example
type TaskMethod struct {
	TaskName   string // taskName
	TaskMethod Method // TaskMethod
}

// Name get the task name
func (t *TaskMethod) Name() string {
	return t.TaskName
}

// Method get the task method
func (t *TaskMethod) Method(ctx *Context) {
	t.TaskMethod(ctx)
}

var _ ITaskMethod = (*TaskMethod)(nil)

// CreateMethod
func CreateMethod(name string, taskMethod func(ctx *Context)) *TaskMethod {
	return &TaskMethod{
		TaskName:   name,
		TaskMethod: taskMethod,
	}
}
