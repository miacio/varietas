package util

import (
	"encoding/json"
	"errors"
	"reflect"
)

// ToJSON
func ToJSON(obj any) string {
	bt, _ := json.Marshal(obj)
	return string(bt)
}

// IsAnyMethod
// check the any have a method
func IsAnyMethod(obj any, method string) bool {
	vo := reflect.ValueOf(obj)
	vt := vo.Type()
	if vo.NumMethod() > 0 {
		for i := 0; i < vo.NumMethod(); i++ {
			if vt.Method(i).Name == method {
				return true
			}
		}
	}
	return false
}

// CallAnyString
func CallAnyString(obj any) (res string, err error) {
	defer func() {
		if r := recover(); r != nil {
			str, ok := r.(string)
			if ok {
				err = errors.New(str)
			} else {
				err = errors.New("call any string panic")
			}
		}
	}()

	method := "String"
	notUse := true

	vo := reflect.ValueOf(obj)
	vt := vo.Type()
	if vo.NumMethod() > 0 {
		for i := 0; i < vo.NumMethod(); i++ {
			if vt.Method(i).Name == method {
				notUse = false
				resVal := vo.Method(i).Call([]reflect.Value{})
				res = resVal[0].String()
			}
		}
	}
	if notUse {
		err = errors.New("unknown String method")
	}

	return res, err
}
