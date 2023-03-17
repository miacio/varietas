# varietas
Experimental business logic variants

# dynamic
Registering an object based on the registration method allows the object to call it through a string

``` go
package dynamic_test

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
	dynamic.Register("test", &tes)

	// params := []reflect.Value{reflect.ValueOf("hello")}
	result, err := dynamic.Call("test", "Goo", nil)
	if err != nil {
		// fmt.Println(result)
		fmt.Printf("call method fail: %v", err)
	}
	fmt.Println(result)

	methods := dynamic.GetMethods("test")
	fmt.Println(methods)
}

func TestV002(t *testing.T) {
	tes := Test{
		Host: "host",
	}
	ctx := dynamic.New()
	ctx.Register("test", &tes)

	// params := []reflect.Value{reflect.ValueOf("hello")}
	result, err := ctx.Call("test", "Goo", nil)
	if err != nil {
		// fmt.Println(result)
		fmt.Printf("call method fail: %v", err)
	}
	fmt.Println(result)

	methods := ctx.GetMethods("test")
	fmt.Println(methods)
}
```
