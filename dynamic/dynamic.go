package dynamic

import (
	"fmt"
	"reflect"
	"sync"
)

// MethodMap
type MethodMap map[string]reflect.Value

// ClassesMethodMap
type ClassesMethodMap map[string]MethodMap

// context
type context struct {
	classesMethodMap ClassesMethodMap
	classesLock      sync.Mutex
}

// New
func New() *context {
	return &context{
		classesMethodMap: make(ClassesMethodMap),
	}
}

// register
func (c *context) register(name string, class any) (int, error) {
	to := reflect.TypeOf(class)
	fmt.Println(to.Kind())
	if to.Kind() != reflect.Struct && to.Kind() != reflect.Pointer {
		return 0, fmt.Errorf("the class type of kind is not a struct")
	}

	if name == "" {
		name = to.String()
	}

	c.classesLock.Lock()
	defer c.classesLock.Unlock()

	if c.classesMethodMap == nil {
		c.classesMethodMap = make(ClassesMethodMap)
	}

	if _, ok := c.classesMethodMap[name]; ok {
		return 0, fmt.Errorf("the current %s class already exists", name)
	}

	methodMap := make(MethodMap)

	vf := reflect.ValueOf(class)
	vft := vf.Type()

	methodNumber := vf.NumMethod()
	if methodNumber > 0 {
		for i := 0; i < methodNumber; i++ {
			methodName := vft.Method(i).Name
			methodMap[methodName] = vf.Method(i)
		}
		c.classesMethodMap[name] = methodMap
		return methodNumber, nil
	}
	return 0, nil
}

// RegisterByName
func (c *context) RegisterByName(name string, class any) (int, error) {
	return c.register(name, class)
}

// RegisterByAny
func (c *context) RegisterByAny(class any) (int, error) {
	return c.register("", class)
}

// call
func (c *context) call(name, method string, params []reflect.Value) ([]reflect.Value, error) {
	methodMap, ok := c.classesMethodMap[name]
	if !ok {
		return nil, fmt.Errorf("the %s class does not exist", name)
	}
	methodValue, ok := methodMap[method]
	if !ok {
		return nil, fmt.Errorf("the %s method does not exist", method)
	}

	result := methodValue.Call(params)
	return result, nil
}

// CallByName
func (c *context) CallByName(name, method string, params []reflect.Value) ([]reflect.Value, error) {
	return c.call(name, method, params)
}

// CallByAny
func (c *context) CallByAny(class any, method string, params []reflect.Value) ([]reflect.Value, error) {
	name := reflect.TypeOf(class).String()
	return c.call(name, method, params)
}

func (c *context) getMethods(name string) []string {
	methodMap, ok := c.classesMethodMap[name]
	if !ok {
		return nil
	}
	methods := make([]string, 0)
	for methodName := range methodMap {
		methods = append(methods, methodName)
	}
	return methods
}

// GetMethods
func (c *context) GetMethodsByName(name string) []string {
	return c.getMethods(name)
}

// GetMethodsByAny
func (c *context) GetMethodsByAny(class any) []string {
	name := reflect.TypeOf(class).String()
	return c.getMethods(name)
}
