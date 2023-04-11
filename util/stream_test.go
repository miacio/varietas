package util_test

import (
	"fmt"
	"testing"

	"github.com/miacio/varietas/util"
)

func TestStream(t *testing.T) {
	v := []int{1, 2, 3, 4, 5, 6}
	sa := util.NewStreamArray(v)

	sa.Skip(1).Limit(10).ForEach(func(v int) int {
		if v > 3 {
			v = 0
		}
		return v
	})
	fmt.Println(sa.Get())
}
