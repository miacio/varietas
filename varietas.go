package varietas

import (
	"errors"
	"reflect"
)

type Context struct {
	classes        map[string]any                      // struct map
	classesMethods map[string]map[string]reflect.Value // classes method names
}

func (c *Context) Register(name string, class any) error {
	if c.classes == nil {
		c.classes = make(map[string]any)
	}
	if c.classesMethods == nil {
		c.classesMethods = make(map[string]map[string]reflect.Value)
	}
	if _, ok := c.classes[name]; ok {
		return errors.New("the current class already exists")
	}
	c.classes[name] = class
	c.classesMethods[name] = make(map[string]reflect.Value)

	vf := reflect.ValueOf(class)
	vfType := vf.Type()

	methodNum := vf.NumMethod()
	for i := 0; i < methodNum; i++ {
		methodName := vfType.Method(i).Name
		c.classesMethods[name][methodName] = vf.Method(i)
	}

	return nil
}

func (c *Context) Get(name string) (any, error) {
	result, ok := c.classes[name]
	if !ok {
		return nil, errors.New("the current class is empty")
	}
	return result, nil
}

func (c *Context) Call(name, method string, params []reflect.Value) ([]reflect.Value, error) {
	_, err := c.Get(name)
	if err != nil {
		return nil, err
	}
	if params == nil {
		return c.classesMethods[name][method].Call([]reflect.Value{}), nil
	}
	result := c.classesMethods[name][method].Call(params)
	if len(result) == 0 {
		return nil, nil
	}
	return result, nil
}
