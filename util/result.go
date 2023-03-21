package util

import (
	"fmt"
	"reflect"
	"strconv"
)

type Result struct {
	ctx any
	err error
}

func Res(ctx any, err error) *Result {
	return &Result{ctx: ctx, err: err}
}

// String
// converts the content of ctx into a string and returns it.
// If a String function exists in ctx,
// the String function will be called and returned
func (rc *Result) String() string {
	ctx := rc.ctx
	ct := reflect.TypeOf(ctx)
	switch ct.Kind() {
	case reflect.String:
		return reflect.ValueOf(ctx).String()
	case reflect.Slice, reflect.Array:
		cv := reflect.New(ct.Elem()).Elem()
		ct = reflect.TypeOf(cv.Interface())
		return fmt.Sprintf("%v", ctx)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(reflect.ValueOf(ctx).Float(), 'f', -1, 64)
	case reflect.Struct, reflect.Ptr:
		val, err := CallAnyString(ctx)
		if err == nil {
			return val
		}
		return fmt.Sprintf("%v", ctx)
	default:
		return fmt.Sprintf("%v", ctx)
	}
}

// Value
func (r *Result) Value() any {
	return r.ctx
}

// ValueTo
// Quickly convert the ctx in the Result to the specified type and return
// If the specified type is inconsistent with the ctx type, it will cause panic
func ValueTo[T any](r *Result, val T) T {
	val = r.Value().(T)
	return val
}

// Error
func (r *Result) Error() error {
	return r.err
}

// HaveError
// return r.err != nil
func (r *Result) HaveError() bool {
	return r.err != nil
}

// OperateError
// Determine whether the error exists.
// If it exists, it will be operated according to the method written by the user,
// thereby reducing the need for developers to repeatedly make if judgments.
func (r *Result) OperateError(v func(error)) {
	if r.err != nil {
		v(r.err)
	}
}
