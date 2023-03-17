package dynamic

import (
	"fmt"
	"reflect"
)

var (
	classesMethodMap ClassesMethodMap
)

// MethodMap
type MethodMap map[string]reflect.Value

// ClassesMethodMap
type ClassesMethodMap map[string]MethodMap

func makeClassesMethodMap() {
	if classesMethodMap == nil {
		classesMethodMap = make(ClassesMethodMap)
	}
}

func init() {
	makeClassesMethodMap()
}

// Register
func Register(name string, class any) (int, error) {
	makeClassesMethodMap()

	if _, ok := classesMethodMap[name]; ok {
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
		classesMethodMap[name] = methodMap
		return methodNumber, nil
	}
	return 0, nil
}

// Call
func Call(name, method string, params []reflect.Value) ([]reflect.Value, error) {
	methodMap, ok := classesMethodMap[name]
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

// GetMethods
func GetMethods(name string) []string {
	methodMap, ok := classesMethodMap[name]
	if !ok {
		return nil
	}
	methods := make([]string, 0)
	for methodName := range methodMap {
		methods = append(methods, methodName)
	}
	return methods
}
