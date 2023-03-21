package util_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/miacio/varietas/util"
)

func TestResult001(t *testing.T) {
	res := util.Res(util.Local().IP())
	res.OperateError(func(err error) {
		log.Fatalf("get ip fail: %v", err)
	})
	// fmt.Println(res.String())
	fmt.Println(util.ValueTo(res, 0))

	var a float64
	a = 123.1
	var b float64
	fmt.Println(util.ValueTo(util.Res(a, nil), b))
}
