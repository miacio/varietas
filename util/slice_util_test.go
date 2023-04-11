package util_test

import (
	"testing"

	"github.com/miacio/varietas/util"
)

func TestSliceMaxNumber(t *testing.T) {
	a := []float32{1, 23, 6, 231, 5, 0, 32}
	max := util.SliceMaxNumber(a)
	t.Fatal(max)
}

func TestSliceMinNumber(t *testing.T) {
	a := []float32{1, 23, 6, 231, 5, 21, 32}
	max := util.SliceMinNumber(a)
	t.Fatal(max)
}
