package dynamic_test

// If use dynamic register not a pointer struct then the struct methods don't use pointer

import (
	"fmt"
	"testing"

	"github.com/miacio/varietas/dynamic"
)

type Test struct {
	Host string
}

func (t *Test) Hello(message string) string {
	fmt.Println("hello")
	return "hello" + message
}

func (t *Test) Goo() {
	fmt.Println("goo")
}

func TestV001(t *testing.T) {
	tes := Test{
		Host: "host",
	}
	ctx := dynamic.New()
	ctx.RegisterByName("test", &tes)

	// params := []reflect.Value{reflect.ValueOf("hello")}
	result, err := ctx.CallByName("test", "Goo", nil)
	if err != nil {
		// fmt.Println(result)
		fmt.Printf("call method fail: %v", err)
	}
	if result != nil {
		fmt.Println(result)
	}

	methods := ctx.GetMethodsByName("test")
	fmt.Println(methods)
}

func TestV002(t *testing.T) {
	tes := Test{
		Host: "host",
	}
	ctx := dynamic.New()
	ctx.RegisterByAny(&tes)

	// params := []reflect.Value{reflect.ValueOf("hello")}
	result, err := ctx.CallByAny(&tes, "Goo", nil)
	if err != nil {
		// fmt.Println(result)
		fmt.Printf("call method fail: %v", err)
	}
	if result != nil {
		fmt.Println(result)
	}

	methods := ctx.GetMethodsByAny(&tes)
	fmt.Println(methods)
}

func TestV003(t *testing.T) {
	tes := Test{
		Host: "host",
	}
	ctx := dynamic.New()
	ctx.RegisterByAny(&tes)

	// params := []reflect.Value{reflect.ValueOf("hello")}

	result, err := ctx.CallByName("*dynamic_test.Test", "Goo", nil)
	if err != nil {
		// fmt.Println(result)
		fmt.Printf("call method fail: %v", err)
	}
	if result != nil {
		fmt.Println(result)
	}

	methods := ctx.GetMethodsByAny(&tes)
	fmt.Println(methods)
}
