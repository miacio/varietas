package varietas_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/miacio/varietas"
)

type Test struct {
}

func (t *Test) Hello(message string) string {
	fmt.Println("hello")
	return "hello" + message
}

func TestV001(t *testing.T) {
	c := varietas.Context{}
	tes := Test{}
	c.Register("test", &tes)

	params := []reflect.Value{reflect.ValueOf("hello")}
	result, err := c.Call("test", "Hello", params)
	if err != nil {
		// fmt.Println(result)
		fmt.Printf("call method fail: %v", err)
	}
	fmt.Println(result)

}
